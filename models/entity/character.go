package entity

import "api-go/enums"

type Character struct {
	ID   int64               `json:"id"`
	Name string              `json:"name"`
	Hp   int32               `json:"hp"`
	Type enums.CharacterType `json:"type"`
}
