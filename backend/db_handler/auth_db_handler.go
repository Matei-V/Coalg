package db_handler

import (
	"database/sql"
	"fmt"
)

type User struct {
	UserName       string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

func Add_User(db *sql.DB, name, password, sessionToken string) {
	stmt := "insert into user(name, parola, sessiontoken) values(?, ?, ?);"
	insert, err := db.Query(stmt, name, password, sessionToken)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func Add_Session_Token(db *sql.DB, name string, token string) {
	stmt := "update user set sessiontoken=? where name=?;"
	alter, err := db.Query(stmt, token, name)
	if err != nil {
		panic(err.Error())
	}
	defer alter.Close()
}

func Get_User(db *sql.DB, name string) (User, bool) {
	var ans User
	ans.UserName = name
	fmt.Print(ans.UserName + "\n")
	stmt := "select parola from user where name='" + name + "';"
	row := db.QueryRow(stmt)
	err := row.Scan(&ans.HashedPassword)
	if ans.HashedPassword == "" {
		return User{}, false
	}
	if err != nil {
		panic(err.Error())
	}
	stmt = "select sessiontoken from user where name='" + name + "';"
	row = db.QueryRow(stmt)
	err = row.Scan(&ans.SessionToken)
	if err != nil {
		panic(err.Error())
	}
	return ans, true
}

func Add_CSRF_Token(db *sql.DB, name string, token string) {
	stmt := "update user set csrftoken=? where name=?;"
	alter, err := db.Query(stmt, token, name)
	if err != nil {
		panic(err.Error())
	}
	defer alter.Close()
}
func Get_User_Password(db *sql.DB, name string) string {
	var ans string
	stmt := "select parola from user where name='admin';"
	row := db.QueryRow(stmt)
	err := row.Scan(&ans)
	if err != nil {
		panic(err.Error())
	}
	return ans
}
