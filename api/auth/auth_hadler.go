package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
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

	if len(username) < 4 {
		In := "Numele de utilizator trebuie să conțină cel puțin 4 caractere."
		var filename = "static/pages/register.html"
		t, err := template.ParseFiles(filename)
		if err != nil {
			panic(err.Error())
		}

		t.Execute(w, In)
		return
	}
	if len(password) < 8 {
		In := "Parola trebuie să conțină cel puțin 8 caractere."
		var filename = "static/pages/register.html"
		t, err := template.ParseFiles(filename)
		if err != nil {
			panic(err.Error())
		}

		t.Execute(w, In)
		return
	}
	_, ok := dbh.Get_User(db, username)
	if ok != false {
		fmt.Println("register get")
		In := "Numele de utilizator este deja utilizat. Alege alt nume."
		var filename = "static/pages/register.html"
		t, err := template.ParseFiles(filename)
		if err != nil {
			panic(err.Error())
		}

		t.Execute(w, In)
		return
	}
	dbh.Add_User(db, username, hashedPassword, "")
	http.Redirect(w, r, "/login-page?registered=true", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Print(username + password + "\n")
	user, ok := dbh.Get_User(db, username)
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		In := "Nume de utilizator sau parolă incorecte."
		var filename = "static/pages/login.html"
		t, err := template.ParseFiles(filename)
		if err != nil {
			panic(err.Error())
		}

		t.Execute(w, In)
		return
	}
	sessionToken := generateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	dbh.Add_Session_Token(db, username, sessionToken)
	http.Redirect(w, r, "/", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if username, err := r.Cookie("username"); err == nil && username != nil {
		dbh.Add_Session_Token(db, username.Value, "")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
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
