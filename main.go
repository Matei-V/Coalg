package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	api "mysql_drive/api"
	dbh "mysql_drive/backend/db_handler"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = dbh.Connect_db()

func exeIndex(w http.ResponseWriter, r *http.Request) {
	var filename = "index.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	err = t.ExecuteTemplate(w, filename, nil)
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

	if !add_cap {
		str := r.PostFormValue("capitol")
		var idx int
		for i := 0; i < len(str); i++ {
			_, err := strconv.Atoi(str[i : i+1])
			if err == nil {
				aux, _ := strconv.Atoi(str[i : i+1])
				idx = idx*10 + aux
			}
		}
		cap = dbh.Create_chap(idx, str[1:])
	}

	dbh.Add_page(db, u, cap, add_cap, pag, true)
}

func add_FileSystem() {
	fs := http.FileServer(http.Dir("template/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	fs = http.FileServer(http.Dir("template/scripts"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", fs))
}

func add_ApiHandlers() {
	http.HandleFunc("/load-menu", api.Load_menu)
	http.HandleFunc("/get-menu", api.Write_menu)
	http.HandleFunc("/add-page", exePageadd)
	http.HandleFunc("/add_new_page", addPage)
	http.HandleFunc("/load-lesson", api.Load_lesson)
	http.HandleFunc("/load-chapters", api.Load_chapters)
}

func main() {
	add_FileSystem()
	add_ApiHandlers()

	http.HandleFunc("/", exeIndex)
	http.ListenAndServe(":8000", nil)
}
