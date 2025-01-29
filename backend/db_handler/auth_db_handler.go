package db_handler

import (
	"database/sql"
	"fmt"
	"strconv"
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
func Get_Lessons_Auth(db *sql.DB, username string) ([]string, []string, []string, []string) {
	var clmn, id_cap, id_less, aprobat []string
	row_num := Last_id(db, "pagina", "id_pag")
	for i := 0; i < row_num; i++ {
		var ans, aux string
		stmt := "select titlu from pagina where id_pag=" + strconv.Itoa(i+1) + " and autor='" + username + "';"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		stmt = "select aprobat from pagina where id_pag=" + strconv.Itoa(i+1) + " and autor='" + username + "';"
		row = db.QueryRow(stmt)
		row.Scan(&aux)
		if ans != "" {
			clmn = append(clmn, ans)
			id_less = append(id_less, strconv.Itoa(i+1))
			aprobat = append(aprobat, aux)
		}
	}
	for i := 0; i < len(clmn); i++ {
		var ans string
		stmt := "select id_unit from capitol where id_cap=(select id_cap from pagina where titlu='" + clmn[i] + "')"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		if ans != "" {
			id_cap = append(id_cap, ans)
		}
	}

	return clmn, id_cap, id_less, aprobat
}

func Get_Lessons_Admin(db *sql.DB, username string) ([]string, []string, []string, []string) {
	var clmn, id_cap, id_less, aprobat []string
	row_num := Last_id(db, "pagina", "id_pag")
	for i := 0; i < row_num; i++ {
		var ans string
		stmt := "select titlu from pagina where id_pag=" + strconv.Itoa(i+1) + " and aprobat=0;"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		if ans != "" {
			clmn = append(clmn, ans)
			id_less = append(id_less, strconv.Itoa(i+1))
			aprobat = append(aprobat, "-1")
		}
	}
	for i := 0; i < len(clmn); i++ {
		var ans string
		stmt := "select id_unit from capitol where id_cap=(select id_cap from pagina where titlu='" + clmn[i] + "')"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		fmt.Print(stmt + "\n")
		fmt.Print(ans + "\n")
		if ans != "" {
			id_cap = append(id_cap, ans)
		}
	}

	return clmn, id_cap, id_less, aprobat
}

func IsAdmin(db *sql.DB, username string) bool {
	var ans string
	stmt := "select admin from user where name='" + username + "';"
	row := db.QueryRow(stmt)
	row.Scan(&ans)

	return ans == "1"
}
