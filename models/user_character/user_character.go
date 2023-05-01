package models

import (
	"api-go/db"
	"api-go/models/entity"
	"database/sql"
	"fmt"
)

func Insert(uc entity.UserCharacter) (id int64, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	query := `INSERT INTO tb_user_character (level, hp_min, slot, user_id, character_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = conn.QueryRow(query, uc.Level, uc.HpMin, uc.Slot, uc.UserId, uc.CharacterId).Scan(&id)
	return
}

func Update(userId int64, uc entity.UserCharacter) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	query := `UPDATE tb_user_character SET level=$1, hp_min=$2, slot=$3 WHERE id=$4 AND user_id=$5`
	res, err := conn.Exec(query, uc.Level, uc.HpMin, uc.Slot, uc.ID, userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Get(id int64, userId int64) (uc entity.UserCharacter, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	query := `SELECT * FROM tb_user_character WHERE id=$1 AND user_id=$2`
	row := conn.QueryRow(query, id, userId)
	err = row.Scan(&uc.ID, &uc.Level, &uc.HpMin, &uc.Slot, &uc.UserId, &uc.CharacterId)
	if err != nil {
		err = fmt.Errorf("nenhum personagem do usuário encontrado")
		return uc, err
	}
	return
}

func GetAll(userId int64) (ucs []entity.UserCharacter, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	query := `SELECT * FROM tb_user_character WHERE user_id=$1`
	rows, err := conn.Query(query, userId)
	if err != nil {
		return
	}
	for rows.Next() {
		var uc entity.UserCharacter
		err = rows.Scan(&uc.ID, &uc.Level, &uc.HpMin, &uc.Slot, &uc.UserId, &uc.CharacterId)
		if err != nil {
			continue
		}
		ucs = append(ucs, uc)
	}
	return
}

func ClearAllSlot(userId int64, ids []int64) (int64, error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return 0, err
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	query := "UPDATE tb_user_character SET slot=$1 WHERE slot >$1 AND id NOT IN ($2, $3, $4) AND user_id=$5"
	s0 := int64(0)
	s1 := int64(0)
	s2 := int64(0)
	if len(ids) > 0 {
		s0 = ids[0]
	}
	if len(ids) > 1 {
		s1 = ids[1]
	}
	if len(ids) > 2 {
		s2 = ids[2]
	}
	res, err := conn.Exec(query, 0, s0, s1, s2, userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
