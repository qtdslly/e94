package cache

import (
	"background/common/logger"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/patrickmn/go-cache"
	"runtime/debug"
	"github.com/mohae/deepcopy"
	"reflect"
)

type Store struct {
	pool     *redis.Pool
	mem      *cache.Cache
	memTtl   int
	redisTTL int
	mutex    *sync.Mutex
	loading  map[string]bool
	lazy     bool
}

type valueWrapper struct {
	Value     string      `json:"value"`
	ValueObj  interface{} `json:"-"`
	ExpiredAt time.Time   `json:"expired_at"`
	CreatedAt time.Time   `json:"created_at"`
}

type StoreLoadFunc func() (interface{}, error)

func NewStore(redisAddr, redisPassword string, redisTTL, memTTL int, lazy bool) *Store {
	s := new(Store)
	s.pool = GetRedisPool(redisAddr, redisPassword)
	s.mem = cache.New(time.Duration(memTTL)*time.Second, time.Second*2)
	s.redisTTL = redisTTL
	s.memTtl = memTTL
	s.mutex = new(sync.Mutex)
	s.loading = make(map[string]bool)
	s.lazy = lazy
	return s
}

func (s *Store) saveValue(key string, data string, ttl int, mem, redis bool) {
	redisTtl := s.redisTTL
	memTtl := s.memTtl
	if ttl > 0 {
		redisTtl = ttl
		memTtl = ttl
	}

	var val valueWrapper
	val.Value = data

	if mem {
		val.ExpiredAt = time.Now().Add(time.Second * time.Duration(memTtl))
		val.CreatedAt = time.Now()
		if s.lazy {
			s.mem.Set(key, &val, time.Duration(memTtl)*time.Second*512)
		} else {
			s.mem.Set(key, &val, time.Duration(memTtl)*time.Second)
		}
	}

	if redis {
		val.ExpiredAt = time.Now().Add(time.Second * time.Duration(redisTtl))
		val.CreatedAt = time.Now()
		bs, _ := json.Marshal(&val)
		if s.lazy {
			err := RedisSetString(key, string(bs), redisTtl*512, s.pool)
			if err != nil {
				logger.Error(err)
			}
		} else {
			err := RedisSetString(key, string(bs), redisTtl, s.pool)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func fWraper(f StoreLoadFunc) (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			logger.Error("panic: ", err)
		}
	}()

	return f()
}

func (s *Store) lazyLoadRedis(key string, ttl int, f StoreLoadFunc) {
	loadingKey := key + "_loading"
	loading, err := RedisGetInt(loadingKey, s.pool)
	if err != nil && err != redis.ErrNil {
		logger.Error(err)
		return
	}
	if loading > 0 {
		return
	}

	RedisSetInt(loadingKey, 1, 60, s.pool)
	defer RedisDelKey(loadingKey, s.pool)

	o, err := fWraper(f)
	if err != nil {
		logger.Error(err)
		return
	}
	data, _ := json.Marshal(o)
	s.saveValue(key, string(data), ttl, true, true)
}

func (s *Store) lazyLoadMem(key string, ttl int, f StoreLoadFunc) {
	s.mutex.Lock()
	_, ok := s.loading[key]
	if ok {
		s.mutex.Unlock()
		return
	}
	s.loading[key] = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.loading, key)
		s.mutex.Unlock()
	}()

	str, err := RedisGetString(key, s.pool)
	if err != nil && err != redis.ErrNil {
		logger.Error(err)
		return
	}
	if str == "" {
		s.lazyLoadRedis(key, ttl, f)
		return
	}

	var val valueWrapper
	err = json.Unmarshal([]byte(str), &val)
	if err != nil {
		logger.Error(err)
		return
	}
	if val.ExpiredAt.UnixNano() < time.Now().UnixNano() && s.lazy {
		// do lazy load
		s.lazyLoadRedis(key, ttl, f)
		return
	}
	s.saveValue(key, val.Value, ttl, true, false)
}

func (s *Store) doCleanLoad(key string, obj interface{}, ttl int, f StoreLoadFunc) error {
	loadingKey := key + "_loading"

	n := time.Now()
	for {
		loading, err := RedisGetInt(loadingKey, s.pool)
		if err != nil && err != redis.ErrNil {
			logger.Error(err)
			return err
		}
		if loading > 0 {
			if time.Now().Unix()-n.Unix() > 10 {
				return errors.New("fetch result for key " + key + " timeout")
			}
			time.Sleep(time.Millisecond * 50)
			continue
		}
		break
	}

	RedisSetInt(loadingKey, 1, 60, s.pool)
	defer RedisDelKey(loadingKey, s.pool)

	o, err := fWraper(f)
	if err != nil {
		logger.Error(err)
		return err
	}
	data, err := json.Marshal(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	s.saveValue(key, string(data), ttl, true, true)
	err = json.Unmarshal(data, obj)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *Store) GetJsonObjectWithExpire(key string, obj interface{}, ttl int, f StoreLoadFunc) error {
	var val valueWrapper
	v, ok := s.mem.Get(key)
	if ok {
		valp := v.(*valueWrapper)
		if valp.ExpiredAt.UnixNano() < time.Now().UnixNano() && s.lazy {
			// do lazy load
			go s.lazyLoadMem(key, ttl, f)
		}
		if valp.ValueObj != nil {
			v := deepcopy.Copy(valp.ValueObj)
			reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(v).Elem())
			return nil
		}
		err := json.Unmarshal([]byte(valp.Value), obj)
		if err != nil {
			logger.Error(err)
			return err
		}
		valp.ValueObj = deepcopy.Copy(obj)
		return nil
	}

	str, err := RedisGetString(key, s.pool)
	if err != nil && err != redis.ErrNil {
		logger.Error(err)
		return err
	}
	if str != "" {
		err := json.Unmarshal([]byte(str), &val)
		if err != nil {
			logger.Error(err)
			return err
		}
		if val.ExpiredAt.UnixNano() < time.Now().UnixNano() && s.lazy {
			// do lazy load
			go s.lazyLoadRedis(key, ttl, f)
		}
		err = json.Unmarshal([]byte(val.Value), obj)
		if err != nil {
			logger.Error(err)
			return err
		}
		s.saveValue(key, val.Value, ttl, true, false)
		return nil
	}

	return s.doCleanLoad(key, obj, ttl, f)
}

func (s *Store) GetJsonObject(key string, obj interface{}, f StoreLoadFunc) error {
	return s.GetJsonObjectWithExpire(key, obj, 0, f)
}

func (s *Store) Delete(keyPattern string) {
}
