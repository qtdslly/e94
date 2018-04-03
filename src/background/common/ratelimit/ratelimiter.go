package ratelimit

import (
	"time"
)

const (
	MaxQPS             = 50
	BucketCapacity     = 500
	MaxWaitTime        = 100 * time.Millisecond
	CacheCleanDuration = 10 * time.Second
)

/* the cache holds the ip and its token bucket */
var cache map[string]*item

/* the timer to maintain the cache, clean it if some ip is no more called for some time */
var cacheCleanTimer *time.Timer

type item struct {
	TokenBucket *Bucket   // token bucket
	AccessedAt  time.Time // last api accessed time
}

func init() {
	cache = make(map[string]*item)
	cacheCleanTimer = time.NewTimer(CacheCleanDuration)
	go cleanCache()
}

/* this method will check if current call has token to access */
func Take(ip string) bool {
	var it *item
	var exist bool
	if it, exist = cache[ip]; !exist {
		now := time.Now()
		tb := NewBucketWithRate(MaxQPS, BucketCapacity)
		it = &item{
			TokenBucket: tb,
			AccessedAt:  now,
		}
		cache[ip] = it
	}
	_, ok := it.TokenBucket.TakeMaxDuration(1, MaxWaitTime)
	return ok
}

func cleanCache() {
	for {
		select {
		case <-cacheCleanTimer.C:
			// clean the cache when the timer is fired
			now := time.Now()
			for ip, item := range cache {
				if now.Sub(item.AccessedAt) >= CacheCleanDuration {
					delete(cache, ip)
				}
			}
			cacheCleanTimer.Reset(CacheCleanDuration)
		}
	}
}
