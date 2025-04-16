package conf

import (
	"KIM/logger"
	"KIM/util"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config Config
type Config struct {
	ServiceID     string
	NodeID        int64
	Listen        string `default:":8080"`
	PublicAddress string
	PublicPort    int `default:"8080"`
	Tags          []string
	ConsulURL     string
	RedisAddrs    string
	Driver        string `default:"mysql"`
	BaseDb        string
	MessageDb     string
	LogLevel      string `default:"INFO"`
}

func (c Config) String() string {
	bts, _ := json.Marshal(c)
	return string(bts)
}

// Init InitConfig
func Init(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/conf")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		logger.Warn(err)
	} else {
		if err := viper.Unmarshal(&config); err != nil {
			return nil, err
		}
	}
	err := envconfig.Process("kim", &config)
	if err != nil {
		return nil, err
	}
	if config.ServiceID == "" {
		localIP := util.GetLocalIP()
		config.ServiceID = fmt.Sprintf("royal_%s", strings.ReplaceAll(localIP, ".", ""))
		arr := strings.Split(localIP, ".")
		if len(arr) == 4 {
			suffix, _ := strconv.Atoi(arr[3])
			config.NodeID = int64(suffix)
		}
	}
	if config.PublicAddress == "" {
		config.PublicAddress = util.GetLocalIP()
	}
	logger.Info(config)
	return &config, nil
}

func InitRedis(addr string, pass string) (*redis.Client, error) {
	redisdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return redisdb, nil
}

// InitFailoverRedis init redis with sentinels
func InitFailoverRedis(masterName string, sentinelAddrs []string, password string, timeout time.Duration) (*redis.Client, error) {
	redisdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: sentinelAddrs,
		Password:      password,
		DialTimeout:   time.Second * 5,
		ReadTimeout:   timeout,
		WriteTimeout:  timeout,
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		logrus.Warn(err)
	}
	return redisdb, nil
}
