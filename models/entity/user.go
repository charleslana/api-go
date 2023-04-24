package entity

type User struct {
	ID          int64            `json:"id"`
	Email       string           `json:"email"`
	Password    string           `json:"password"`
	Name        *string          `json:"name"`
	Permissions []UserPermission `json:"permissions"`
}
