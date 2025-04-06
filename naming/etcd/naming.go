package etcd

import (
	"KIM/communication"
	"KIM/naming"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

type EtcdNaming struct {
	cli      *clientv3.Client
	watches  map[string]context.CancelFunc
	lock     sync.RWMutex
	rootPath string
}

const defaultRoot = "/kim-registry"

func NewEtcdNaming(endpoints []string) (naming.Naming, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdNaming{
		cli:      cli,
		watches:  make(map[string]context.CancelFunc),
		rootPath: defaultRoot,
	}, nil
}

func (e *EtcdNaming) key(serviceName, serviceID string) string {
	return fmt.Sprintf("%s/%s/%s", e.rootPath, serviceName, serviceID)
}

func (e *EtcdNaming) Register(service communication.ServiceRegistration) error {
	key := e.key(service.ServiceName(), service.ServiceID())
	val, err := json.Marshal(service)
	if err != nil {
		return err
	}
	_, err = e.cli.Put(context.Background(), key, string(val))
	return err
}

func (e *EtcdNaming) Deregister(serviceID string) error {
	ctx := context.Background()
	// find and delete by ID (match any serviceName)
	resp, err := e.cli.Get(ctx, e.rootPath, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		var reg naming.DefaultServiceRegistration
		_ = json.Unmarshal(kv.Value, &reg)
		if reg.Id == serviceID {
			_, err := e.cli.Delete(ctx, string(kv.Key))
			return err
		}
	}
	return nil
}

func (e *EtcdNaming) Find(serviceName string, tags ...string) ([]communication.ServiceRegistration, error) {
	ctx := context.Background()
	prefix := fmt.Sprintf("%s/%s/", e.rootPath, serviceName)
	resp, err := e.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var services []communication.ServiceRegistration
	for _, kv := range resp.Kvs {
		var reg naming.DefaultServiceRegistration
		if err := json.Unmarshal(kv.Value, &reg); err != nil {
			continue
		}
		if matchTags(reg.Tags, tags) {
			services = append(services, reg)
		}
	}
	return services, nil
}

func matchTags(serviceTags []string, queryTags []string) bool {
	if len(queryTags) == 0 {
		return true
	}
	tagSet := make(map[string]struct{})
	for _, t := range serviceTags {
		tagSet[t] = struct{}{}
	}
	for _, q := range queryTags {
		if _, ok := tagSet[q]; !ok {
			return false
		}
	}
	return true
}

func (e *EtcdNaming) Subscribe(serviceName string, callback func([]communication.ServiceRegistration)) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if _, ok := e.watches[serviceName]; ok {
		return errors.New("already watching service")
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.watches[serviceName] = cancel

	go func() {
		watchChan := e.cli.Watch(ctx, fmt.Sprintf("%s/%s/", e.rootPath, serviceName), clientv3.WithPrefix())
		for range watchChan {
			services, err := e.Find(serviceName)
			if err == nil && callback != nil {
				callback(services)
			}
		}
	}()
	return nil
}

func (e *EtcdNaming) Unsubscribe(serviceName string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	cancel, ok := e.watches[serviceName]
	if ok {
		cancel()
		delete(e.watches, serviceName)
	}
	return nil
}
