package project

import (
	"errors"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (repo *ProjectRepository) Create(project *Project) *Project {
	if err := repo.db.Create(&project).Error; err != nil {
		panic(err)
	}

	return project
}

func (repo *ProjectRepository) Save(project *Project) *Project {
	if err := repo.db.Save(&project).Error; err != nil {
		panic(err)
	}

	return project
}

func (repo *ProjectRepository) FindByID(userID, projectID uint64) *Project {
	var ret Project
	q := repo.db.Where("user_id = ?", userID).Where("id = ?", projectID)
	if err := q.Take(&ret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(err)
	}

	return &ret
}

func (repo *ProjectRepository) FindAll(userID uint64) []Project {
	var ret []Project
	if err := repo.db.Where("user_id = ?", userID).Find(&ret).Error; err != nil {
		panic(err)
	}

	return ret
}
