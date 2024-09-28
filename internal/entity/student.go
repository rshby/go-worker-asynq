package entity

import (
	"context"
	"time"
)

type Student struct {
	ID             uint64    `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement" json:"id"`
	IdentityNumber string    `gorm:"column:identity_number;type:varchar(256);not null" json:"identity_number"`
	Name           string    `gorm:"column:name;type:varchar(256);default:null" json:"name"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;not null;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;not null;autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (s *Student) TableName() string {
	return "students"
}

// abstraction
type (
	StudentRepository interface {
		InsertStudent(ctx context.Context, input *Student) (*Student, error)
		GetStudentByID(ctx context.Context, id uint64) (*Student, error)
		GetStudentByIdentityNumber(ctx context.Context, identityNumber string) (*Student, error)
		LockInsertStudentByIdentityNumber(ctx context.Context, identityNumber string) (func(), error)
	}

	StudentService interface {
		InsertStudent(ctx context.Context, request *RequestInsertStudent) error
		InsertStudentBulk(ctx context.Context, request *RequestInsertStudentBulk) error
	}
)

// dto
type (
	RequestInsertStudent struct {
		IdentityNumber string `json:"identity_number" validate:"required" example:"123456"`
		Name           string `json:"name,omitempty" example:"John Doe"`
	}

	RequestInsertStudentBulk struct {
		Data []RequestInsertStudent `json:"data" validate:"required"`
	}
)

// Validate is function to validate request insert student bulk
func (r *RequestInsertStudentBulk) Validate() error {
	return validate.Struct(*r)
}
