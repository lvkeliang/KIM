package database

import (
	"KIM/protocol/protoImpl"
	"KIM/services"
	"fmt"
	"github.com/go-redis/redis/v7"
	"google.golang.org/protobuf/proto"
	"time"
)

const LocationExpired = time.Hour * 48

type RedisStorage struct {
	cli *redis.Client
}

func NewRedisStorage(cli *redis.Client) services.SessionStorage {
	return &RedisStorage{
		cli: cli,
	}
}

func (r *RedisStorage) Add(session *protoImpl.Session) error {
	loc := services.Location{
		ChannelId: session.ChannelId,
		GateId:    session.GateId,
	}

	// 先保存Location
	locKey := KeyLocation(session.Account, "")
	err := r.cli.Set(locKey, loc.Bytes(), LocationExpired).Err()
	if err != nil {
		return err
	}

	// 再保存Session
	snKey := KeySession(session.ChannelId)
	buf, _ := proto.Marshal(session)
	err = r.cli.Set(snKey, buf, LocationExpired).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisStorage) Delete(account string, channelId string) error {
	locKey := KeyLocation(account, "")
	err := r.cli.Del(locKey).Err()
	if err != nil {
		return err
	}

	snKey := KeySession(channelId)
	err = r.cli.Del(snKey).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStorage) Get(channelId string) (*protoImpl.Session, error) {
	snKey := KeySession(channelId)
	bts, err := r.cli.Get(snKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, services.ErrSessionNil
		}
		return nil, err
	}
	var session protoImpl.Session
	_ = proto.Unmarshal(bts, &session)
	return &session, nil
}

// GetLocations 批量获取位置信息,主要用于群聊时读取群成员定位信息
func (r *RedisStorage) GetLocations(accounts ...string) ([]*services.Location, error) {
	keys := KeyLocations(accounts...)
	list, err := r.cli.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}
	var result = make([]*services.Location, 0)
	for _, l := range list {
		if l == nil {
			continue
		}
		var loc services.Location
		_ = loc.Unmarshal([]byte(l.(string)))
		result = append(result, &loc)
	}
	if len(result) == 0 {
		return nil, services.ErrSessionNil
	}
	return result, nil
}

// GetLocation 获取某个账号的信息，device作为预留参数，用于支持多设备登录
func (r *RedisStorage) GetLocation(account string, device string) (*services.Location, error) {
	key := KeyLocation(account, device)
	bts, err := r.cli.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, services.ErrSessionNil
		}
		return nil, err
	}
	var loc services.Location
	_ = loc.Unmarshal(bts)
	return &loc, nil
}

func KeySession(channel string) string {
	return fmt.Sprintf("login:sn:%s", channel)
}

func KeyLocation(account, device string) string {
	if device == "" {
		return fmt.Sprintf("login:loc:%s", account)
	}
	return fmt.Sprintf("login:loc:%s:%s", account, device)
}

func KeyLocations(accounts ...string) []string {
	arr := make([]string, len(accounts))
	for i, account := range accounts {
		arr[i] = KeyLocation(account, "")
	}
	return arr
}
