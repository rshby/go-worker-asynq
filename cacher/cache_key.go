package cacher

import "fmt"

func GetStudentByIDCacheKey(id uint64) string {
	return createCacheKey(fmt.Sprintf("cache:object:students:id:%d", id))
}

func GetStudentByIdentityNumberCacheKey(identityNumber string) string {
	return createCacheKey(fmt.Sprintf("cache:object:students:identity_number:%s", identityNumber))
}
