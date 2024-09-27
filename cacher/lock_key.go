package cacher

import "fmt"

func InsertStudentByIdentityNumberLockKey(identityNumber string) string {
	return createCacheKey(fmt.Sprintf("lock:insert_student:identity_number:%s", identityNumber))
}
