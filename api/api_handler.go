package api

import (
	"database/sql"
	"fmt"
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
		less := dbh.Get_Lessons(db, i+1)
		for j := 0; j < len(less); j++ {
			fmt.Printf(less[j] + " " + caps[i] + "\n")
			htmlStr += "<form hx-post='/load-lesson' hx-target='#content' hx-swap='innerHTML'>"
			htmlStr += "<input type='hidden' name='id' value='" + less[j] + "'>"
			htmlStr += "<input type='submit' name='submit' value='" + less[j] + "'>"
			htmlStr += "</form>"
		}
		htmlStr += "</ul></li>"
	}
	htmlStr += "</ul>"

	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}

func Load_lesson(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("id")
	desc := dbh.Get_pag(db, title)
	htmlStr := "<h3>" + title + "</h3>"
	htmlStr += desc
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}
