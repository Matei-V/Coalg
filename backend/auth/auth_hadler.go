package auth

import (
	"database/sql"
	"fmt"
	dbh "mysql_drive/backend/db_handler"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName       string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var db *sql.DB = dbh.Connect_db()

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	fmt.Print(username + " " + password + "\n")
	hashedPassword, _ := hashPassword(password)

	dbh.Add_User(db, username, hashedPassword)
}
func Login(w http.ResponseWriter, r *http.Request) {}

func Logout(w http.ResponseWriter, r *http.Request) {}

func Protected(w http.ResponseWriter, r *http.Request) {}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
