package cache

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

// 获取redis连接池对象
func GetRedisPool(address, password string) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     50,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dial("tcp", address, password)
		},
	}
	return pool
}

func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

func RedisTest(address, password string) error {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer c.Close()
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return err
		}
	}
	_, err = c.Do("PING")
	return err
}

// 自增
func RedisIncr(key string, ttl int, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	_, err := conn.Do("INCR", key)
	if err != nil {
		return err
	}
	// ttl>0时认为需要设置过期时间
	if ttl > 0 {
		_, err = conn.Do("expire", key, ttl)
		return err
	}
	return nil
}

// 获取所有key
func RedisGetKeysByPattern(pattern string, pool *redis.Pool) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return []string{}, err
	}

	data, err := redis.Strings(conn.Do("keys", pattern))
	return data, err
}

func RedisSetUint64(key string, value uint64, ttl int, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}
	_, err := conn.Do("SETEX", key, ttl, value)
	return err

}

func RedisIncrBy(key string, incrAmount uint64, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	_, err := conn.Do("INCRBY", key, incrAmount)
	if err != nil {
		return err
	}
	return nil
}

func RedisSetKeyExpire(key string, ttl int, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}
	_, err := conn.Do("EXPIRE", key, ttl)
	if err != nil {
		return err
	}
	return nil
}

func RedisGetUint64(key string, pool *redis.Pool) (uint64, error) {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return 0, err
	}

	n, err := redis.Uint64(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}
	return n, err
}

func RedisSetObject(key string, obj interface{}, ttl int, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	if err := conn.Err(); err != nil {
		return err
	}
	_, err = conn.Do("SETEX", key, ttl, string(bs))
	return err
}

func RedisGetObject(key string, obj interface{}, pool *redis.Pool) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return false, err
	}

	s, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(s), obj)
	if err != nil {
		return false, err
	}
	return true, nil
}

func RedisSetString(key, value string, ttl int, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}
	_, err := conn.Do("SETEX", key, ttl, value)
	return err
}

func RedisGetString(key string, pool *redis.Pool) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return "", err
	}

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return s, nil
}

func RedisSetInt(key string, value, ttl int, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}
	_, err := conn.Do("SETEX", key, ttl, value)
	return err

}

func RedisGetInt(key string, pool *redis.Pool) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return 0, err
	}

	n, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}
	return n, err
}

func RedisDelKey(key string, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	var err error
	if err = conn.Err(); err != nil {
		return err
	}

	_, err = conn.Do("DEL", key)
	return err
}
