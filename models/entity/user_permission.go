package entity

type UserPermission struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	UserId int64  `json:"userId"`
}
