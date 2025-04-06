package etcd

import (
	"KIM/communication"
	"KIM/logger"
	"KIM/naming"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
)

const (
	defaultRootPath = "/kim-registry"
	defaultLeaseTTL = 15 * time.Second
)

type EtcdNaming struct {
	client   *clientv3.Client
	leaseMap sync.Map // serviceID -> leaseID
	watchers map[string]context.CancelFunc
	rootPath string
	mu       sync.RWMutex
}

// etcdServiceEntry 统一存储结构
type etcdServiceEntry struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Protocol string            `json:"protocol"`
	Tags     []string          `json:"tags"`
	Meta     map[string]string `json:"meta"`
	LastSeen time.Time         `json:"last_seen"`
}

func NewEtcdNaming(endpoints []string) (naming.Naming, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdNaming{
		client:   client,
		watchers: make(map[string]context.CancelFunc),
		rootPath: defaultRootPath,
	}, nil
}

func (n *EtcdNaming) buildKey(serviceName, serviceID string) string {
	return fmt.Sprintf("%s/%s/%s", n.rootPath, serviceName, serviceID)
}

// Register 通过接口方法获取字段
func (n *EtcdNaming) Register(service communication.ServiceRegistration) error {
	entry := etcdServiceEntry{
		ID:       service.ServiceID(),
		Name:     service.ServiceName(),
		Address:  service.PublicAddress(),
		Port:     service.PublicPort(),
		Protocol: service.GetProtocol(),
		Tags:     service.GetTags(),
		Meta:     service.GetMeta(),
		LastSeen: time.Now(),
	}

	// 创建租约
	leaseResp, err := n.client.Grant(context.Background(), int64(defaultLeaseTTL.Seconds()))
	if err != nil {
		return fmt.Errorf("create lease failed: %v", err)
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal service entry failed: %v", err)
	}

	// 存储数据并绑定租约
	key := n.buildKey(service.ServiceName(), service.ServiceID())
	if _, err := n.client.Put(context.Background(), key, string(data), clientv3.WithLease(leaseResp.ID)); err != nil {
		return fmt.Errorf("put service entry failed: %v", err)
	}

	// 保持租约存活
	n.leaseMap.Store(service.ServiceID(), leaseResp.ID)
	go n.keepAlive(service.ServiceID())

	return nil
}

func (n *EtcdNaming) keepAlive(serviceID string) {
	ticker := time.NewTicker(defaultLeaseTTL / 2)
	defer ticker.Stop()

	for range ticker.C {
		leaseID, ok := n.leaseMap.Load(serviceID)
		if !ok {
			return
		}

		// 1. 执行真实的租约续期
		if _, err := n.client.KeepAliveOnce(context.Background(), leaseID.(clientv3.LeaseID)); err != nil {
			logger.Warnf("Lease keepalive failed for %s: %v", serviceID, err)
			return
		}

		// 2. 更新 LastSeen（可选）
		key := n.buildKey("", serviceID)
		resp, err := n.client.Get(context.Background(), key)
		if err != nil || len(resp.Kvs) == 0 {
			continue
		}

		var entry etcdServiceEntry
		if json.Unmarshal(resp.Kvs[0].Value, &entry) == nil {
			entry.LastSeen = time.Now()
			data, _ := json.Marshal(entry)
			n.client.Put(context.Background(), key, string(data), clientv3.WithLease(leaseID.(clientv3.LeaseID)))
		}
	}
}

func (n *EtcdNaming) Deregister(serviceID string) error {
	if leaseID, ok := n.leaseMap.Load(serviceID); ok {
		n.client.Revoke(context.Background(), leaseID.(clientv3.LeaseID))
		n.leaseMap.Delete(serviceID)
	}

	_, err := n.client.Delete(context.Background(), n.buildKey("", serviceID))
	return err
}

func (n *EtcdNaming) Find(serviceName string, tags ...string) ([]communication.ServiceRegistration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := n.client.Get(ctx,
		n.buildKey(serviceName, ""),
		clientv3.WithPrefix(),
	)
	if err != nil {
		return nil, fmt.Errorf("etcd get failed: %v", err)
	}

	var services []communication.ServiceRegistration
	now := time.Now()

	for _, kv := range resp.Kvs {
		var entry etcdServiceEntry
		if err := json.Unmarshal(kv.Value, &entry); err != nil {
			continue
		}

		// 实施健康检查
		if now.Sub(entry.LastSeen) > defaultLeaseTTL*2 {
			continue
		}

		// 转换为统一响应结构
		service := &naming.DefaultServiceRegistration{
			Id:       entry.ID,
			Name:     entry.Name,
			Address:  entry.Address,
			Port:     entry.Port,
			Protocol: entry.Protocol,
			Tags:     entry.Tags,
			Meta:     entry.Meta,
		}

		if matchTags(service.Tags, tags) {
			services = append(services, service)
		}
	}

	return services, nil
}

func (n *EtcdNaming) Subscribe(serviceName string, callback func([]communication.ServiceRegistration)) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.watchers[serviceName]; ok {
		return errors.New("service already watched")
	}

	ctx, cancel := context.WithCancel(context.Background())
	n.watchers[serviceName] = cancel

	watchChan := n.client.Watch(ctx,
		n.buildKey(serviceName, ""),
		clientv3.WithPrefix(),
	)

	go func() {
		// 正确处理Watch通道
		for watchResp := range watchChan {
			// 1. 处理可能的错误
			if watchResp.Err() != nil {
				logger.Warnf("Watch error: %v", watchResp.Err())
				return
			}

			// 2. 触发全量更新（保持与原逻辑一致）
			if services, err := n.Find(serviceName); err == nil {
				callback(services)
			}
		}
	}()

	return nil
}

func (n *EtcdNaming) Unsubscribe(serviceName string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if cancel, ok := n.watchers[serviceName]; ok {
		cancel()
		delete(n.watchers, serviceName)
	}
	return nil
}

func matchTags(serviceTags []string, requiredTags []string) bool {
	tagSet := make(map[string]struct{})
	for _, tag := range serviceTags {
		tagSet[tag] = struct{}{}
	}

	for _, reqTag := range requiredTags {
		if _, exists := tagSet[reqTag]; !exists {
			return false
		}
	}
	return true
}
