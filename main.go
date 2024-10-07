package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type unit struct {
	id_unit int
}

type chap struct {
	id_cap int
	titlu  string
}

type pag struct {
	id_pag int
	titlu  string
}

type res struct {
	id_res int
	titlu  string
	link   string
	limba  string
}

type prob struct {
	id_pb int
	titlu string
	link  string
	limba string
}

func add_page(db *sql.DB, u unit, c chap, p pag, r []res, pb []prob) {
	insert, err := db.Query("insert into pagina(titlu, id_cap) values(?, ?)", p.titlu, c.id_cap)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < len(r); i++ {
		insert, err = db.Query("insert into resursa(titlu, link, limba,  id_pag) values(?, ?, ?, ?)",
			r[i].titlu, r[i].link, r[i].limba, p.id_pag)
		if err != nil {
			panic(err.Error())
		}
	}
	for i := 0; i < len(pb); i++ {
		insert, err = db.Query("insert into problema(titlu, link, limba,  id_pag) values(?, ?, ?, ?)",
			pb[i].titlu, pb[i].link, pb[i].limba, p.id_pag)
		if err != nil {
			panic(err.Error())
		}
	}
	defer insert.Close()
}

func main() {
	pass := os.Getenv("MYSQL_PASS")
	db, err := sql.Open("mysql", "root:"+pass+"@tcp(localhost:3306)/coalg_db")
	if err != nil {
		fmt.Println("error")
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err.Error())
	}

	u := unit{9}
	c := chap{1, "Introducere in limbajul C++"}
	p := pag{1, "Complexitate asimptotica"}

	r := []res{{1, "PbInfo", "pbinfo.ro", "Romana"}, {2, "InfoArena", "infoarena.ro", "Romana"}}
	pb := []prob{{1, "PbInfo", "pbinfo.ro", "Romana"}, {2, "Kilonova", "kilo.ro", "Romana"}}
	add_page(db, u, c, p, r, pb)
}
