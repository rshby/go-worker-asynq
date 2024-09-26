package entity

import "context"

type TaskQueue interface {
	Enqueue(ctx context.Context, taskName string, data any) error
	Stop()
}
