package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	api "mysql_drive/api"
	dbh "mysql_drive/backend/db_handler"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = dbh.Connect_db()

func exeIndex(w http.ResponseWriter, r *http.Request) {
	my_url, _ := url.Parse(r.URL.String())
	param, _ := url.ParseQuery(my_url.RawQuery)

	unit := param.Get("clasa")

	var filename = "index.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	err = t.ExecuteTemplate(w, filename, unit)
	if err != nil {
		panic(err.Error())
	}
}

func exePageadd(w http.ResponseWriter, r *http.Request) {
	var filename = "add_page.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	err = t.ExecuteTemplate(w, filename, nil)
	if err != nil {
		panic(err.Error())
	}
}

func addPage(w http.ResponseWriter, r *http.Request) {
	id_unitate, _ := strconv.Atoi(r.PostFormValue("unitate"))
	u := dbh.Create_unit(id_unitate)

	add_cap, _ := strconv.ParseBool(r.PostFormValue("add-cap"))
	cap := dbh.Create_chap(dbh.Last_id(db, "capitol", "id_cap")+1, r.PostFormValue("capitol"))
	pag := dbh.Create_pag(dbh.Last_id(db, "pagina", "id_pag")+1, r.PostFormValue("titlu"), r.PostFormValue("descriere"))

	nr_res, _ := strconv.Atoi(r.PostFormValue("nr_res"))
	var titlu_res, link_res []string
	var id_res []int
	init_id := dbh.Last_id(db, "resursa", "id_res")
	for i := 0; i < nr_res; i++ {
		titlu_res = append(titlu_res, r.PostFormValue("res_"+strconv.Itoa(i+1)))
		link_res = append(link_res, r.PostFormValue("link_res_"+strconv.Itoa(i+1)))
		id_res = append(id_res, init_id+i+1)
	}
	res := dbh.Create_res(id_res, titlu_res, link_res)
	add_res := (nr_res != 0)

	nr_pb, _ := strconv.Atoi(r.PostFormValue("nr_pb"))
	var titlu_pb, link_pb []string
	var id_pb []int
	init_id = dbh.Last_id(db, "problema", "id_prob")
	for i := 0; i < nr_pb; i++ {
		titlu_pb = append(titlu_pb, r.PostFormValue("pb_"+strconv.Itoa(i+1)))
		link_pb = append(link_pb, r.PostFormValue("link_pb_"+strconv.Itoa(i+1)))
		id_pb = append(id_pb, init_id+i+1)
	}
	pb := dbh.Create_pb(id_pb, titlu_pb, link_pb)
	add_pb := (nr_pb != 0)

	if !add_cap {
		str := r.PostFormValue("capitol")
		var idx, i int
		aux, err := strconv.Atoi(str[0:1])
		for i = 0; i < len(str) && err == nil; i++ {
			aux, err = strconv.Atoi(str[i : i+1])
			if err == nil {
				idx = idx*10 + aux
			}
		}
		cap = dbh.Create_chap(idx, str[i:])
	}

	dbh.Add_page(db, u, cap, add_cap, pag, true, res, add_res, pb, add_pb)
}

func add_ApiHandlers() {
	http.HandleFunc("/load-menu", api.Load_menu)
	http.HandleFunc("/get-menu", api.Write_menu)
	http.HandleFunc("/add-page", exePageadd)
	http.HandleFunc("/add_new_page", addPage)
	http.HandleFunc("/load-lesson", api.Load_lesson)
	http.HandleFunc("/load-chapters", api.Load_chapters)
	http.HandleFunc("/add_new_chapter", api.Add_new_chapter)
	http.HandleFunc("/add_to_existing_chapter", api.Add_to_existing_chapter)
	http.HandleFunc("/add_res", api.Add_res)
	http.HandleFunc("/update_nr_res", api.Update_nr_res)
	http.HandleFunc("/add_pb", api.Add_pb)
	http.HandleFunc("/update_nr_pb", api.Update_nr_pb)
}

func main() {
	add_ApiHandlers()
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))

	http.HandleFunc("/", exeIndex)
	http.ListenAndServe(":8000", nil)
}
