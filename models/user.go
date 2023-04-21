package models

import (
	"api-go/db"
	"api-go/models/entity"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

func Insert(user entity.User) (id int64, err error) {
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
	query := `INSERT INTO tb_user (email, password) VALUES ($1, $2) RETURNING id`
	password, _ := HashPassword(user.Password)
	err = conn.QueryRow(query, user.Email, password).Scan(&id)
	return
}

func Update(id int64, name string) (int64, error) {
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
	query := `UPDATE tb_user SET name=$1 WHERE id=$2`
	res, err := conn.Exec(query, name, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func Get(id int64) (user entity.User, err error) {
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
	query := `SELECT * FROM tb_user WHERE id=$1`
	row := conn.QueryRow(query, id)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Name)
	return
}

func GetAll() (users []entity.User, err error) {
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
	query := `SELECT * FROM tb_user`
	rows, err := conn.Query(query)
	if err != nil {
		return
	}
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return
}

func Delete(id int64) (int64, error) {
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
	query := `DELETE FROM tb_user WHERE id=$1`
	res, err := conn.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func HashPassword(password string) (pass string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	pass = string(bytes)
	return
}

func CheckPassword(password string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
