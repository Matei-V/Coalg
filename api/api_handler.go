package api

import (
	"database/sql"
	"html/template"
	"net/http"

	dbh "mysql_drive/backend/db_handler"
)

var db *sql.DB = dbh.Connect_db()

func Load_menu(w http.ResponseWriter, r *http.Request) {
	htmlStr := "template/menu_tmp.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, nil)
}

func Write_menu(w http.ResponseWriter, r *http.Request) {
	caps := dbh.Get_Chapters(db)
	htmlStr := "<ul class='chapters'>\n"
	for i := 0; i < len(caps); i++ {
		htmlStr += "<li>" + caps[i] + "<ul>"
		less := dbh.Get_Lessons(db, i)
		for j := 0; j < len(less); j++ {
			htmlStr += "<li>" + less[j] + "</li>"
		}
		htmlStr += "</ul></li>"
	}
	htmlStr += "</ul>"

	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}
