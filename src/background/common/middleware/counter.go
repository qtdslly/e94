package middleware

import (
	"background/common/cache"
	"background/common/logger"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

var mutex *sync.Mutex
var requestMap map[request]uint32
var successMap map[request]uint32
var failMap map[request]uint32

const (
	requestAliveTTL = 3600 * 24 * 3 // 过期时间

	apiTypeReq     = `req`
	apiTypeSuccess = `success`
	apiTypeFail    = `fail`
)

type request struct {
	Path string
	Time int64
}

func init() {
	mutex = new(sync.Mutex)
	requestMap = make(map[request]uint32)
	successMap = make(map[request]uint32)
	failMap = make(map[request]uint32)
}

// 获取redis key
func getApiCacheKey(typ string, timeNow int64, url string) string {
	hour := timeNow / 3600
	return fmt.Sprintf(`api-%s-%s-%d-%d`, typ, url, hour, timeNow)
}

func flushCounter(pool *redis.Pool) {
	keyMap := make(map[string]uint32)
	n := time.Now().Unix()
	mutex.Lock()
	for k, v := range requestMap {
		cacheKey := getApiCacheKey(apiTypeReq, k.Time, k.Path)
		keyMap[cacheKey] = v
		if n - k.Time > 5 {
			delete(requestMap, k)
		} else {
			requestMap[k] = 0
		}
	}
	for k, v := range successMap {
		cacheKey := getApiCacheKey(apiTypeSuccess, k.Time, k.Path)
		keyMap[cacheKey] = v
		if n-k.Time > 5 {
			delete(successMap, k)
		} else {
			successMap[k] = 0
		}
	}
	for k, v := range failMap {
		cacheKey := getApiCacheKey(apiTypeFail, k.Time, k.Path)
		keyMap[cacheKey] = v
		if n-k.Time > 5 {
			delete(failMap, k)
		} else {
			failMap[k] = 0
		}
	}
	mutex.Unlock()

	for k, v := range keyMap {
		var err error
		if v == 0 {
			continue
		}
		if err = cache.RedisIncrBy(k, uint64(v), pool); err != nil {
			logger.Error(err)
		}
		if err = cache.RedisSetKeyExpire(k, requestAliveTTL, pool); err != nil {
			logger.Error(err)
		}
	}
}

func incCounter(r request, m map[request]uint32) {
	mutex.Lock()
	defer mutex.Unlock()

	// increment the request count
	if v, exists := m[r]; !exists {
		m[r] = 1
	} else {
		m[r] = v + 1
	}
}

func CounterHandler(pool *redis.Pool) gin.HandlerFunc {
	go func() {
		for {
			flushCounter(pool)
			time.Sleep(time.Millisecond * 500)
		}
	}()

	return func(c *gin.Context) {
		r := request{
			Path: c.Request.URL.Path,
			Time: time.Now().Unix(),
		}
		incCounter(r, requestMap)
		c.Next()
		if c.Writer.Status() == http.StatusOK {
			incCounter(r, successMap)
		} else {
			incCounter(r, failMap)
		}
	}
}
