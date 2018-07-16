package service

import (
	"background/common/cache"
	"fmt"
)

var cacheStore *cache.Store

func InitCache(redisAddr, redisPassword string) {
	cacheStore = cache.NewStore(redisAddr, redisPassword, 60, 10, true)
}

func GetCacheKey(entity string, appId, versionId, contentType, contentId uint32, args ...interface{}) string {
	key := fmt.Sprintf("%s_app_id_%d_version_id_%d_content_type_%d_content_id_%d", entity, appId, versionId, contentType, contentId)
	for _, k := range args {
		key += fmt.Sprint(k)
	}
	return key
}

func GetCacheObject(key string, obj interface{}, f cache.StoreLoadFunc) error {
	return cacheStore.GetJsonObject(key, obj, f)
}

func GetCacheObjectWithExpire(key string, obj interface{}, ttl int, f cache.StoreLoadFunc) error {
	return cacheStore.GetJsonObjectWithExpire(key, obj, ttl, f)
}

// 更新缓存
func SetCacheObjectWithExpire(key string, obj interface{}, ttl int) error {
	return cacheStore.SetJsonObjectWithExpire(key, obj, ttl)
}