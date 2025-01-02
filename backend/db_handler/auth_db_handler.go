package db_handler

import (
	"database/sql"
)

func Add_User(db *sql.DB, name string, password string) {
	stmt := "insert into user(name, parola) values(?, ?);"
	insert, err := db.Query(stmt, name, password)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}
