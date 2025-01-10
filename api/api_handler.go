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
	htmlStr := "static/template/menu_tmp.html"
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
			htmlStr += "<a href='/unit?clasa=" + unit + "&lectie=" + less[j] + "'>" + less[j]
			htmlStr += "</a><br>"
		}
	}

	tmpl, _ := template.New("t").Parse(htmlStr)

	tmpl.Execute(w, nil)
}

func Load_chapters(w http.ResponseWriter, r *http.Request) {
	type In struct {
		Caps, Id string
	}

	var Input []In

	aux, aux_id := dbh.Get_Chapters(db, "9")
	for i := 0; i < len(aux); i++ {
		Input = append(Input, In{Caps: aux[i], Id: aux_id[i]})
	}

	aux, aux_id = dbh.Get_Chapters(db, "10")
	for i := 0; i < len(aux); i++ {
		Input = append(Input, In{Caps: aux[i], Id: aux_id[i]})
	}

	aux, aux_id = dbh.Get_Chapters(db, "11")
	for i := 0; i < len(aux); i++ {
		Input = append(Input, In{Caps: aux[i], Id: aux_id[i]})
	}
	htmlStr := "static/template/select_chapters.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, Input)
}

func Add_new_chapter(w http.ResponseWriter, r *http.Request) {
	htmlStr := "static/template/new_chapter.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, nil)
}

func Add_to_existing_chapter(w http.ResponseWriter, r *http.Request) {
	htmlStr := "static/template/ex_chapter.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, nil)
}
func Add_res(w http.ResponseWriter, r *http.Request) {
	nr_res, _ := strconv.Atoi(r.PostFormValue("nr_res"))
	nr_res++
	htmlStr := "static/template/add_res.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, strconv.Itoa(nr_res))
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
	htmlStr := "static/template/add_pb.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, strconv.Itoa(nr_pb))
}

func Update_nr_pb(w http.ResponseWriter, r *http.Request) {
	nr_pb, _ := strconv.Atoi(r.PostFormValue("nr_pb"))
	htmlStr := "<input class='nr_pb' type='hidden' name='nr_pb' value='" + strconv.Itoa(nr_pb) + "'>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}
func Load_lesson(w http.ResponseWriter, r *http.Request) {
	my_url, _ := url.Parse(r.URL.String())
	param, _ := url.ParseQuery(my_url.RawQuery)
	title := param.Get("lectie")

	if title == "" {
		htmlStr := "<div id='content' style='position:relative; left:240px'>"
		htmlStr += "</div>"

		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
		return
	}

	desc := dbh.Get_pag(db, title)
	titlu_res, link_res := dbh.Get_res(db, title)
	htmlStr := "<div id='content' style='position:relative; left:240px'>"
	htmlStr += "<h3>" + title + "</h3>"
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

	htmlStr += "</div>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}
