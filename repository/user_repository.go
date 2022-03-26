package repository

import (
	"go-tech-blog/model"
)

func UserCreate(user *model.User) (*model.User, error) {
	tx := db.Begin()

	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return user, nil
}
