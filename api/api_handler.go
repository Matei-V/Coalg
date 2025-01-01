package api

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	dbh "mysql_drive/backend/db_handler"
)

var db *sql.DB = dbh.Connect_db()

func Load_menu(w http.ResponseWriter, r *http.Request) {
	my_url, _ := url.Parse(r.URL.String())
	param, _ := url.ParseQuery(my_url.RawQuery)
	unit := param.Get("clasa")

	htmlStr := "template/menu_tmp.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, unit)
}

func Write_menu(w http.ResponseWriter, r *http.Request) {
	my_url, _ := url.Parse(r.URL.String())
	param, _ := url.ParseQuery(my_url.RawQuery)

	unit := param.Get("clasa")
	caps, _ := dbh.Get_Chapters(db, unit)
	var htmlStr string
	for i := 0; i < len(caps); i++ {
		if caps[i] == "" {
			continue
		}
		htmlStr += "<h3 class='capitol'>" + caps[i] + "</h3>"
		less := dbh.Get_Lessons(db, i+1)
		for j := 0; j < len(less); j++ {
			htmlStr += "<form hx-post='/load-lesson' hx-target='#content' hx-swap='innerHTML'>"
			htmlStr += "<input  type='hidden' name='id' value='" + less[j] + "'>"
			htmlStr += "<input class='lesson' type='submit' name='submit' value='" + less[j] + "'>"
			htmlStr += "</form>"
		}
	}

	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}

func Load_chapters(w http.ResponseWriter, r *http.Request) {
	htmlStr := "<select class='text-input' name='capitol'>"
	caps, id := dbh.Get_Chapters(db, "9")
	aux, _ := dbh.Get_Chapters(db, "10")
	caps = append(caps, aux...)
	aux, _ = dbh.Get_Chapters(db, "11")
	caps = append(caps, aux...)
	for i := 0; i < len(caps); i++ {
		htmlStr += "<option value='" + id[i] + caps[i] + "' >" + caps[i] + "</option>"
	}

	htmlStr += "</select>"
	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}

func Add_new_chapter(w http.ResponseWriter, r *http.Request) {
	htmlStr := "<input class='text-input' type='text' name='capitol' placeholder='Capitol'>\n"
	htmlStr += "<button class='button' hx-post='/add_to_existing_chapter' hx-target='.Chapter' hx-swap='innerHTML'>Add to exsisting Chapter</button>\n"
	htmlStr += "<input type='hidden' name='add-cap' value='true'>"

	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}

func Add_to_existing_chapter(w http.ResponseWriter, r *http.Request) {
	htmlStr := "<div hx-post='/load-chapters' hx-swap='outerHTML' hx-trigger='load'></div>\n"
	htmlStr += "<button class='button' hx-post='/add_new_chapter' hx-target='.Chapter' hx-swap='innerHTML'>Add new Chapter</button>"
	htmlStr += "<input type='hidden' name='add-cap' value='false'>"

	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Add_res(w http.ResponseWriter, r *http.Request) {
	nr_res, _ := strconv.Atoi(r.PostFormValue("nr_res"))
	nr_res++
	htmlStr := "<input class='text-input' type='text' placeholder='Nume Resursa' name='res_" + strconv.Itoa(nr_res) + "' >\n"
	htmlStr += "<input class='text-input' type='text' placeholder='Link Resursa' name='link_res_" + strconv.Itoa(nr_res) + "' ><br>\n"

	htmlStr += "<input class='nr_res' type='hidden' hx-post='/update_nr_res' hx-trigger='load' hx-swap='outerHTML' name='nr_res' value='" + strconv.Itoa(nr_res) + "'>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Update_nr_res(w http.ResponseWriter, r *http.Request) {
	nr_res, _ := strconv.Atoi(r.PostFormValue("nr_res"))
	htmlStr := "<input class='nr_res' type='hidden' name='nr_res' value='" + strconv.Itoa(nr_res) + "'>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Add_pb(w http.ResponseWriter, r *http.Request) {
	nr_pb, _ := strconv.Atoi(r.PostFormValue("nr_pb"))
	nr_pb++
	htmlStr := "<input class='text-input' type='text' placeholder='Nume Problema' name='pb_" + strconv.Itoa(nr_pb) + "' >\n"
	htmlStr += "<input class='text-input' type='text' placeholder='Link Problema' name='link_pb_" + strconv.Itoa(nr_pb) + "' ><br>\n"

	htmlStr += "<input class='nr_pb' type='hidden' hx-post='/update_nr_pb' hx-trigger='load' hx-swap='outerHTML' name='nr_pb' value='" + strconv.Itoa(nr_pb) + "'>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Update_nr_pb(w http.ResponseWriter, r *http.Request) {
	nr_pb, _ := strconv.Atoi(r.PostFormValue("nr_pb"))
	htmlStr := "<input class='nr_pb' type='hidden' name='nr_pb' value='" + strconv.Itoa(nr_pb) + "'>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Load_lesson(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("id")
	desc := dbh.Get_pag(db, title)

	titlu_res, link_res := dbh.Get_res(db, title)
	htmlStr := "<h3>" + title + "</h3>"
	htmlStr += desc
	htmlStr += "<br>"

	htmlStr += "<h3>Resursa</h3>"
	for i := 0; i < len(titlu_res); i++ {
		htmlStr += "<a href='" + link_res[i] + "'>" + titlu_res[i] + "</a><br>"
	}

	titlu_pb, link_pb := dbh.Get_pb(db, title)
	htmlStr += "<h3>Problema</h3>"
	for i := 0; i < len(titlu_pb); i++ {
		htmlStr += "<a href='" + link_pb[i] + "'>" + titlu_pb[i] + "</a><br>"
	}
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}
