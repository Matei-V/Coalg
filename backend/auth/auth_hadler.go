package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	dbh "mysql_drive/backend/db_handler"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB = dbh.Connect_db()

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	hashedPassword, _ := hashPassword(password)

	dbh.Add_User(db, username, hashedPassword, "")

	fmt.Fprintln(w, "Register succesful")
	http.Redirect(w, r, "/", 301)
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Print(username + password + "\n")
	user, ok := dbh.Get_User(db, username)
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		fmt.Fprintln(w, "login unsuccesful")
		return
	}

	sessionToken := generateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	dbh.Add_Session_Token(db, username, sessionToken)
	http.Redirect(w, r, "/", 301)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	username, _ := r.Cookie("username")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	dbh.Add_Session_Token(db, username.Value, "")
	http.Redirect(w, r, "/", 302)
}

func Authorize(r *http.Request) error {
	username, err := r.Cookie("username")
	if err != nil {
		return errors.New("user not found")
	}
	user, ok := dbh.Get_User(db, username.Value)
	if !ok {
		return errors.New("user not found")
	}

	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		return errors.New("Unauthorized")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(lenght int) string {
	bytes := make([]byte, lenght)

	if _, err := rand.Read(bytes); err != nil {
		fmt.Print("failed to generate token")
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
