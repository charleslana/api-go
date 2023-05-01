package models

import (
	"api-go/db"
	"api-go/models/entity"
	"database/sql"
	"fmt"
)

func Get(id int64) (c entity.Character, err error) {
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
	query := `SELECT * FROM tb_character WHERE id=$1`
	row := conn.QueryRow(query, id)
	err = row.Scan(&c.ID, &c.Name, &c.Hp, &c.Type)
	if err != nil {
		err = fmt.Errorf("nenhum personagem encontrado")
		return c, err
	}
	return
}
