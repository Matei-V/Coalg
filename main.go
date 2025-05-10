package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	api "mysql_drive/api"
	"mysql_drive/api/auth"
	dbh "mysql_drive/backend/db_handler"

	_ "github.com/go-sql-driver/mysql"
)

type Input3 struct {
	Aux1, Aux2, Aux3 string
}

var db *sql.DB = dbh.Connect_db()

func exeIndex(w http.ResponseWriter, r *http.Request) {
	var filename = "index.html"
	t := template.Must(template.ParseFiles(filename))
	err := t.ExecuteTemplate(w, filename, nil)
	if err != nil {
		panic(err.Error())
	}

}

func exePageadd(w http.ResponseWriter, r *http.Request) {
	var filename = "static/pages/add_page.html"
	if err := auth.Authorize(r); err == nil {
		t, err := template.ParseFiles(filename)
		if err != nil {
			panic(err.Error())
		}

		t.Execute(w, nil)
		return
	}

	filename = "static/pages/unauthorized.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	t.Execute(w, nil)
	return

}

func exeRegister(w http.ResponseWriter, r *http.Request) {
	var filename = "static/pages/register.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	t.Execute(w, nil)
}

func exeUnit(w http.ResponseWriter, r *http.Request) {
	my_url, _ := url.Parse(r.URL.String())
	param, _ := url.ParseQuery(my_url.RawQuery)

	unit := param.Get("clasa")
	less := param.Get("id")

	var in Input3
	in.Aux1 = unit
	in.Aux2 = unit
	in.Aux3 = less

	var filename = "static/pages/unit.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	t.Execute(w, in)
}

func exeLogin(w http.ResponseWriter, r *http.Request) {
	var filename = "static/pages/login.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	t.Execute(w, nil)
}
func exeAcc(w http.ResponseWriter, r *http.Request) {
	var filename = "static/pages/account.html"
	t, err := template.ParseFiles(filename)
	if err != nil {
		panic(err.Error())
	}

	t.Execute(w, nil)
}
func addPage(w http.ResponseWriter, r *http.Request) {
	id_unitate, _ := strconv.Atoi(r.PostFormValue("unitate"))
	autor, _ := r.Cookie("username")
	u := dbh.Create_unit(id_unitate)

	add_cap, _ := strconv.ParseBool(r.PostFormValue("add-cap"))
	cap := dbh.Create_chap(dbh.Last_id(db, "capitol", "id_cap")+1, r.PostFormValue("capitol"))
	pag := dbh.Create_pag(dbh.Last_id(db, "pagina", "id_pag")+1, r.PostFormValue("titlu"), r.PostFormValue("descriere"))

	nr_res, _ := strconv.Atoi(r.PostFormValue("nr_res"))
	var titlu_res, link_res []string
	var id_res []int
	init_id := dbh.Last_id(db, "resursa", "id_res")
	idx := 0
	for i := 0; i < nr_res; i++ {
		if r.PostFormValue("res_"+strconv.Itoa(i+1)) == "" {
			continue
		}
		titlu_res = append(titlu_res, r.PostFormValue("res_"+strconv.Itoa(i+1)))
		link_res = append(link_res, r.PostFormValue("link_res_"+strconv.Itoa(i+1)))
		id_res = append(id_res, init_id+idx+1)
		idx++
	}
	res := dbh.Create_res(id_res, titlu_res, link_res)
	add_res := (nr_res != 0)

	nr_pb, _ := strconv.Atoi(r.PostFormValue("nr_pb"))
	var titlu_pb, link_pb []string
	var id_pb []int
	init_id = dbh.Last_id(db, "problema", "id_prob")
	idx = 0
	for i := 0; i < nr_pb; i++ {
		if r.PostFormValue("pb_"+strconv.Itoa(i+1)) == "" {
			continue
		}
		titlu_pb = append(titlu_pb, r.PostFormValue("pb_"+strconv.Itoa(i+1)))
		link_pb = append(link_pb, r.PostFormValue("link_pb_"+strconv.Itoa(i+1)))
		id_pb = append(id_pb, init_id+idx+1)
		idx++
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

	dbh.Add_page(db, u, cap, add_cap, pag, true, res, add_res, pb, add_pb, autor.Value)
}

func add_ApiHandlers() {
	http.HandleFunc("/load-menu", api.Load_menu)
	http.HandleFunc("/get-menu", api.Write_menu)
	http.HandleFunc("/add-page", exePageadd)
	http.HandleFunc("/register-page", exeRegister)
	http.HandleFunc("/login-page", exeLogin)
	http.HandleFunc("/add_new_page", addPage)
	http.HandleFunc("/load-lesson", api.Load_lesson)
	http.HandleFunc("/load-chapters", api.Load_chapters)
	http.HandleFunc("/add_new_chapter", api.Add_new_chapter)
	http.HandleFunc("/add_to_existing_chapter", api.Add_to_existing_chapter)
	http.HandleFunc("/add_res", api.Add_res)
	http.HandleFunc("/update_nr_res", api.Update_nr_res)
	http.HandleFunc("/add_pb", api.Add_pb)
	http.HandleFunc("/update_nr_pb", api.Update_nr_pb)
	http.HandleFunc("/register", auth.Register)
	http.HandleFunc("/login", auth.Login)
	http.HandleFunc("/logout", auth.Logout)
	http.HandleFunc("/auth-cont", api.Auth_cont)
	http.HandleFunc("/unit", exeUnit)
	http.HandleFunc("/my-less", api.MyAcc)
	http.HandleFunc("/account", exeAcc)
	http.HandleFunc("/aproved", api.Aproba)
	http.HandleFunc("/reject", api.Respinge)
}

func main() {
	add_ApiHandlers()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", exeIndex)
	http.ListenAndServe(":8000", nil)
}
