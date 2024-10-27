package main

import (
	"database/sql"
	"html/template"
	"net/http"

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
}

func main() {
	//add_FileSystem()
	add_ApiHandlers()

	http.HandleFunc("/", exeIndex)
	http.ListenAndServe(":8080", nil)
}
