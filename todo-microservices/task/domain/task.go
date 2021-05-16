package task

import (
	"time"
	pb "todo/proto/task"
)

type StatusType int32

type Task struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string
	Status    pb.Status
	ProjectID uint64
	UserID    uint64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func NewTask(name string, userID, projectID uint64) *Task {
	return &Task{
		Name:      name,
		Status:    pb.Status_WAITING,
		UserID:    userID,
		ProjectID: projectID,
	}
}
