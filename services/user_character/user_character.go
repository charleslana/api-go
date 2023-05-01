package services

import (
	"api-go/models/entity"
	models2 "api-go/models/user_character"
)

func Create(uc entity.UserCharacter) (id int64, err error) {
	id, err = models2.Insert(uc)
	return id, err
}

func Update(userId int64, uc entity.UserCharacter) (rows int64, err error) {
	rows, err = models2.Update(userId, uc)
	return rows, err
}

func Get(id int64, userId int64) (uc entity.UserCharacter, err error) {
	uc, err = models2.Get(id, userId)
	return uc, err
}

func List(userId int64) (ucs []entity.UserCharacter, err error) {
	ucs, err = models2.GetAll(userId)
	return ucs, err
}

func CalculateHp(uc entity.UserCharacter) (hp int32) {
	hp = uc.Character.Hp * uc.Level
	return hp
}

func ClearAllSlot(userId int64, ids []int64) (rows int64, err error) {
	rows, err = models2.ClearAllSlot(userId, ids)
	return rows, err
}
