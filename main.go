package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:marynela04@tcp(localhost:3306)/coalg_db")
	if err != nil {
		fmt.Println("error")
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		//panic(err.Error())
		fmt.Println("nu mere vericule")
	}

	insert, err := db.Query("insert into unit values(2);")
	if err != nil {

	}
	defer insert.Close()
}
