package cacher

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/utils"
)

func FindFromCacheByKey[T any](cache CacheManager, key string) (item T, mutex *redsync.Mutex, err error) {
	var cachedData any

	cachedData, mutex, err = cache.GetOrLock(key)
	if err != nil || cachedData == nil {
		return
	}

	cachedDataByte, _ := cachedData.([]byte)
	if cachedDataByte == nil {
		return
	}

	if err = utils.JSONUnmarshal(cachedDataByte, &item); err != nil {
		return
	}

	return
}

func FindFromCacheByKeyWithoutMutex(cache CacheManager, cacheKey string) (string, error) {
	cachedData, err := cache.Get(cacheKey)
	if err != nil {
		return "", err
	}

	bt, _ := cachedData.([]byte)
	return string(bt), nil
}

// SafeUnlock safely unlock mutex
func SafeUnlock(mutex *redsync.Mutex) {
	if mutex != nil {
		_, _ = mutex.Unlock()
	}
}

func StoreNil(ctx context.Context, cache CacheManager, cacheKey string) {
	if err := cache.StoreNil(cacheKey); err != nil {
		logrus.WithContext(ctx).WithField("cacheKey", cacheKey).Error(err)
	}
}

func FindHashMemberFromBucketAndCacheKey[T any](cache CacheManager, bucket, key string) (item T, mu *redsync.Mutex, err error) {
	reply, mu, err := cache.GetHashMemberOrLock(bucket, key)
	if err != nil {
		return
	}

	if reply == nil {
		return
	}

	bt, ok := reply.([]byte)
	if !ok {
		err = ErrFailedCastMultiResponse
		logrus.WithField("reply", reply).Error(err)
		return
	}

	if err = utils.JSONUnmarshal(bt, &item); err != nil {
		return
	}
	return
}

func FindMultiResponseFromCacheByKey(cache CacheManager, bucket, key string) (multiResponse *MultiResponse, mu *redsync.Mutex, err error) {
	reply, mu, err := cache.GetHashMemberOrLock(bucket, key)
	if err != nil {
		return
	}

	if reply == nil {
		return
	}

	bt, ok := reply.([]byte)
	if !ok {
		err = ErrFailedCastMultiResponse
		logrus.WithField("reply", reply).Error(err)
		return nil, nil, err
	}

	multiResponse, err = NewMultiResponseFromByte(bt)
	return
}
