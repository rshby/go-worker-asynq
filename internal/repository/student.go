package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/cacher"
	"go-worker-asynq/config"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/utils"
	"gorm.io/gorm"
)

type studentRepository struct {
	db    *gorm.DB
	cache cacher.CacheManager
}

// NewStudentRepository is method to create instance studentRepository
func NewStudentRepository(db *gorm.DB, cache cacher.CacheManager) entity.StudentRepository {
	return &studentRepository{
		db:    db,
		cache: cache,
	}
}

// InsertStudent is function to insert data student
func (s *studentRepository) InsertStudent(ctx context.Context, input *entity.Student) (*entity.Student, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"input":   utils.Dump(input),
	})

	// create to database
	if err := s.db.WithContext(ctx).Model(&entity.Student{}).Create(input).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	// success insert
	return input, nil
}

// GetStudentByID is method to get data student by id
func (s *studentRepository) GetStudentByID(ctx context.Context, id uint64) (*entity.Student, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"id":      id,
	})

	cacheKey := cacher.GetStudentByIDCacheKey(id)
	if config.EnableCache() {
		cachedItem, mutex, err := cacher.FindFromCacheByKey[*entity.Student](s.cache, cacheKey)
		defer cacher.SafeUnlock(mutex)
		if err != nil {
			logger.Error(err)
		}

		if mutex == nil {
			if cachedItem != nil {
				logger.WithField("cacheKey", cacheKey).Infof("returning data student from redis cache")
				return cachedItem, nil
			}
		}
	}

	// get from database mysql
	var student entity.Student
	err := s.db.WithContext(ctx).Model(&entity.Student{}).Take(&student, id).Error
	switch err {
	case nil:
		if config.EnableCache() {
			if err = s.cache.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.Dump(&student))); err != nil {
				logger.Error(err)
			}
		}

		return &student, nil
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}

// GetStudentByIdentityNumber is method to get data student by identityNumber
func (s *studentRepository) GetStudentByIdentityNumber(ctx context.Context, identityNumber string) (*entity.Student, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context":        utils.DumpIncomingContext(ctx),
		"identityNumber": identityNumber,
	})

	cacheKey := cacher.GetStudentByIdentityNumberCacheKey(identityNumber)
	if config.EnableCache() {
		cachedItem, mutex, err := cacher.FindFromCacheByKey[uint64](s.cache, cacheKey)
		defer cacher.SafeUnlock(mutex)
		if err != nil {
			logger.Error(err)
		}

		if mutex == nil {
			if cachedItem > 0 {
				return s.GetStudentByID(ctx, cachedItem)
			}
		}
	}

	// get data from db
	var id uint64
	if err := s.db.WithContext(ctx).Model(&entity.Student{}).Where("id = ?", id).Pluck("id", &id).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	// if not found
	if id == 0 {
		return nil, nil
	}

	if config.EnableCache() {
		if err := s.cache.StoreWithoutBlocking(cacher.NewItem(cacheKey, id)); err != nil {
			logger.Error(err)
		}
	}

	return s.GetStudentByID(ctx, id)
}

// LockInsertStudentByIdentityNumber is method to lock process insert student by identityNumber
func (s *studentRepository) LockInsertStudentByIdentityNumber(ctx context.Context, identityNumber string) (func(), error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context":        utils.DumpIncomingContext(ctx),
		"identityNumber": identityNumber,
	})

	lockKey := cacher.InsertStudentByIdentityNumberLockKey(identityNumber)
	mutex, err := s.cache.AcquireLock(lockKey)
	if err != nil {
		logger.Error(err)
	}

	return func() {
		cacher.SafeUnlock(mutex)
	}, err
}
