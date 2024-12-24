package api

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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

func Load_chapters(w http.ResponseWriter, r *http.Request) {
	htmlStr := "<select name='capitol'>"
	caps := dbh.Get_Chapters(db)

	for i := 0; i < len(caps); i++ {
		htmlStr += "<option value='" + strconv.Itoa(i+1) + caps[i] + "' >" + caps[i] + "</option>"
	}

	htmlStr += "</option>"
	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}

func Add_new_chapter(w http.ResponseWriter, r *http.Request) {
	htmlStr := "<input type='text' name='capitol' placeholder='Capitol' id='input'>\n"
	htmlStr += "<button hx-post='/add_to_existing_chapter' hx-target='.Chapter' hx-swap='innerHTML'>Add to exsisting Chapter</button>\n"
	htmlStr += "<input type='hidden' name='add-cap' value='true'>"

	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}

func Add_to_existing_chapter(w http.ResponseWriter, r *http.Request) {
	htmlStr := "<div hx-post='/load-chapters' hx-swap='outerHTML' hx-trigger='load'></div>\n"
	htmlStr += "<button hx-post='/add_new_chapter' hx-target='.Chapter' hx-swap='innerHTML'>Add new Chapter</button>"
	htmlStr += "<input type='hidden' name='add-cap' value='false'>"

	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Add_res(w http.ResponseWriter, r *http.Request) {
	nr_res, _ := strconv.Atoi(r.PostFormValue("nr_res"))
	nr_res++
	htmlStr := "<input type='text' placeholder='Nume Resursa' name='res_" + strconv.Itoa(nr_res) + "' ><br>\n"

	htmlStr += "<input class='nr_res' type='hidden' hx-post='/update_nr_res' hx-trigger='load' hx-swap='outerHTML' name='nr_res' value='" + strconv.Itoa(nr_res) + "'>"
	fmt.Print(nr_res)
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Update_nr_res(w http.ResponseWriter, r *http.Request) {
	nr_res, _ := strconv.Atoi(r.PostFormValue("nr_res"))
	nr_res++
	htmlStr := "<input class='nr_res' type='hidden' name='nr_res' value='" + strconv.Itoa(nr_res) + "'>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Add_pb(w http.ResponseWriter, r *http.Request) {
	nr_pb, _ := strconv.Atoi(r.PostFormValue("nr_pb"))
	nr_pb++
	htmlStr := "<input type='text' placeholder='Nume Problema' name='pb_" + strconv.Itoa(nr_pb) + "' ><br>\n"

	htmlStr += "<input class='nr_pb' type='hidden' hx-post='/update_nr_pb' hx-trigger='load' hx-swap='outerHTML' name='nr_pb' value='" + strconv.Itoa(nr_pb) + "'>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Update_nr_pb(w http.ResponseWriter, r *http.Request) {
	nr_pb, _ := strconv.Atoi(r.PostFormValue("nr_pb"))
	nr_pb++
	htmlStr := "<input class='nr_pb' type='hidden' name='nr_pb' value='" + strconv.Itoa(nr_pb) + "'>"
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
