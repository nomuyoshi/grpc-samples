package user

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) *User {
	if err := r.db.Create(user).Error; err != nil {
		panic(err)
	}

	return user
}

func (r *UserRepository) FindByID(userID uint64) *User {
	var u User
	if err := r.db.Where("id = ?", userID).Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(err)
	}

	return &u
}

func (r *UserRepository) FindByEmail(email string) *User {
	var u User
	if err := r.db.Where("email = ?", email).Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(err)
	}

	return &u
}
