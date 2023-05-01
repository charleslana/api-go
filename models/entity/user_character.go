package entity

type UserCharacter struct {
	ID          int64     `json:"id"`
	Level       int32     `json:"level"`
	HpMin       int32     `json:"hpMin"`
	HpMax       int32     `json:"hpMax"`
	Slot        int8      `json:"slot"`
	Character   Character `json:"character"`
	UserId      int64     `json:"userId"`
	CharacterId int64     `json:"characterId"`
}
