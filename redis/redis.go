package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/navi-tt/cache"
	"math/rand"
	"strings"
	"time"
)

var Nil = redis.Nil

var (
	invalidAddress = errors.New("redis addr invalid")
)

type RedisCache struct {
	redis.Cmdable
	client          redisClient
	Master          string // 如果使用哨兵模式
	Addr            []string
	Password        string
	DB              int
	MaxRetries      int // 重试次数
	MaxRetryBackoff int // 最大重试间隔，单位ms

}

type redisClient interface {
	redis.Cmdable

	Close() error
}

func NewRedisCache() cache.Cache {
	return &RedisCache{}
}

func (r *RedisCache) Init(conf string) error {
	if conf == "" {
		return cache.InvalidConfig
	}

	var cf map[string]interface{}
	json.Unmarshal([]byte(conf), &cf)

	r.Master = cf["master"].(string)
	if val, ok := cf["master"].(string); !ok {
		r.Master = ""
	} else {
		r.Master = val
	}

	if val, ok := cf["addr"].(string); !ok {
		return invalidAddress
	} else {
		addr := strings.Split(val, "-")
		r.Addr = addr
	}

	if val, ok := cf["password"].(string); !ok {
		r.Password = ""
	} else {
		r.Password = val
	}

	if val, ok := cf["db"].(int); !ok {
		r.DB = 0
	} else {
		r.DB = val
	}

	if val, ok := cf["master"].(int); !ok {
		r.MaxRetries = 0
	} else {
		if val >= 3 {
			val = 3
		}
		r.MaxRetries = val
	}

	if val, ok := cf["maxRetryBackoff"].(int); !ok {
		r.MaxRetryBackoff = 0
	} else {
		r.MaxRetryBackoff = val
	}

	if r.Master != "" {
		opts := redis.FailoverOptions{
			MasterName:      r.Master,
			SentinelAddrs:   r.Addr,
			Password:        r.Password,
			DB:              r.DB,
			MaxRetries:      r.MaxRetries,
			MaxRetryBackoff: time.Duration(r.MaxRetryBackoff) * time.Millisecond,
		}
		r.client = redis.NewFailoverClient(&opts)
	} else {
		if len(r.Addr) == 1 {
			opts := redis.Options{
				Addr:            r.Addr[0],
				Password:        r.Password,
				DB:              r.DB,
				MaxRetries:      r.MaxRetries,                                        // 最大重试次数
				MaxRetryBackoff: time.Duration(r.MaxRetryBackoff) * time.Millisecond, // 最大重试间隔
			}
			r.client = redis.NewClient(&opts)
		} else {
			opts := redis.ClusterOptions{
				Addrs:           r.Addr,
				Password:        r.Password,
				MaxRetries:      r.MaxRetries,                                        // 最大重试次数
				MaxRetryBackoff: time.Duration(r.MaxRetryBackoff) * time.Millisecond, // 最大重试间隔
			}
			r.client = redis.NewClusterClient(&opts)
		}
	}

	if err := r.client.Ping().Err(); err != nil {
		r.client = nil
		return err
	}

	r.Cmdable = r.client

	return nil
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	fmt.Printf("[freecache Get key:{%s}]\n", key)

	if strings.EqualFold(key, "") {
		return nil, cache.InvalidKey
	}

	ret := r.client.Get(key)
	if err := ret.Err(); err != nil {
		if err == Nil {
			return nil, cache.ErrNotFound
		}
		return nil, err
	}

	val := ret.Val()

	return val, nil
}

func (r *RedisCache) Set(key string, val interface{}, expireTime int) error {
	fmt.Printf("[freecache set key:{%s} value:{%v} timeout:{%d}]\n", key, val, expireTime)

	if strings.EqualFold(key, "") {
		return cache.InvalidKey
	}

	err := r.client.Set(key, val, time.Duration(randExpire(expireTime))*time.Second)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (r *RedisCache) Delete(key string) error {
	fmt.Printf("[freecache del key:{%s}]\n", key)

	if strings.EqualFold(key, "") {
		return cache.InvalidKey
	}

	cmd := r.client.Del(key)
	return cmd.Err()
}

func (r *RedisCache) IsExist(key string) bool {
	ret := r.client.Get(key)

	if err := ret.Err(); err != nil {
		if err == Nil {
			return false
		}
		return false
	}

	return true
}

func init() {
	cache.Register(cache.REDIS_CACHE, NewRedisCache)
}

func randExpire(base int) int {
	return base + rand.Intn(base/3)
}
