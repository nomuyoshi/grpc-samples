package project

import "time"

type Project struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string
	UserID    uint64
	CreatedAt *time.Time
}

func NewProject(name string, userID uint64) *Project {
	return &Project{
		Name:   name,
		UserID: userID,
	}
}
