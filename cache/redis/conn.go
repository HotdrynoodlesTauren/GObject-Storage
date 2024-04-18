package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
	redisHost = "127.0.0.1:6397"
	// redisPass = "testupload" // 不确定需不需要
)

// newRedisPool: create a redis pool
func newRedisPool() *redis.Pool {
	return &redis.Pool {
		MaxIdle: 50,
		MaxActive: 30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// establish a connection
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// authenticate the connection
			// if _, err = c.Do("AUTH", redisPass); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, nil
		},
		TestOnBorrow: func(coon redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := coon.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}