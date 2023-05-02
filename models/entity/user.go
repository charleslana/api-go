package entity

type User struct {
	ID          int64            `json:"id"`
	Email       string           `json:"email"`
	Password    string           `json:"password"` //`json:"-"`
	Name        *string          `json:"name"`
	Permissions []UserPermission `json:"permissions"`
	Characters  []UserCharacter  `json:"characters"`
}
