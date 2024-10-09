package db_handler

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
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

func Add_page(db *sql.DB, u unit, c chap, insert_chap bool, p pag, insert_pag bool, r []res, pb []prob) {
	if insert_chap {
		insert, err := db.Query("insert into capitol(titlu, id_unit) values(?, ?)", p.titlu, u.id_unit)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	if insert_pag {
		insert, err := db.Query("insert into pagina(titlu, id_cap) values(?, ?)", p.titlu, c.id_cap)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	for i := 0; i < len(r); i++ {
		insert, err := db.Query("insert into resursa(titlu, link, limba,  id_pag) values(?, ?, ?, ?)",
			r[i].titlu, r[i].link, r[i].limba, p.id_pag)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	for i := 0; i < len(pb); i++ {
		insert, err := db.Query("insert into problema(titlu, link, limba,  id_pag) values(?, ?, ?, ?)",
			pb[i].titlu, pb[i].link, pb[i].limba, p.id_pag)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
}

func Last_id(db *sql.DB, table string, id_type string) int {
	var ans int
	stmt := "select " + id_type + " from " + table + " order by " + id_type + " desc limit 1"
	row := db.QueryRow(stmt)
	err := row.Scan(&ans)
	if err != nil {
		panic(err.Error())
	}
	return ans
}

func Connect_db() *sql.DB {
	pass := os.Getenv("MYSQL_PASS")
	db, err := sql.Open("mysql", "root:"+pass+"@tcp(localhost:3306)/coalg_db")
	if err != nil {
		fmt.Println("error")
	}

	err = db.Ping()

	if err != nil {
		panic(err.Error())
	}
	return db
}

func Get_Chapters(db *sql.DB) []string {
	var clmn []string
	row_num := Last_id(db, "capitol", "id_cap")

	for i := 0; i < row_num; i++ {
		var ans string
		stmt := "select titlu from capitol where id_cap=" + strconv.Itoa(i+1) + ";"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		clmn = append(clmn, ans)
	}
	return clmn
}
