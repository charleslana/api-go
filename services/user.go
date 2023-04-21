package services

import (
	"api-go/models"
	"api-go/models/entity"
	"fmt"
)

func Create(user entity.User) (id int64, err error) {
	if user.Email == "" {
		err = fmt.Errorf("email em branco")
		return 0, err
	}
	if user.Password == "" {
		err = fmt.Errorf("senha em branco")
		return 0, err
	}
	id, err = models.Insert(user)
	return id, err
}
