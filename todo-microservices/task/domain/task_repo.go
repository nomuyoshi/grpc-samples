package task

import (
	"errors"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Save(task *Task) *Task {
	if err := r.db.Save(task).Error; err != nil {
		panic(err)
	}

	return task
}

func (r *TaskRepository) FindByID(id, userID uint64) *Task {
	var ret Task
	if err := r.db.Where("id = ?", id).Where("user_id = ?", userID).Take(&ret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(err)
	}

	return &ret
}

func (r *TaskRepository) FindByUserID(userID uint64) []Task {
	var ret []Task
	if err := r.db.Where("user_id = ?", userID).Find(&ret).Error; err != nil {
		panic(err)
	}

	return ret
}

func (r *TaskRepository) FindByUserIDAndProjectID(userID, projectID uint64) []Task {
	var ret []Task
	if err := r.db.
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID).
		Find(&ret).
		Error; err != nil {
		panic(err)
	}

	return ret
}
