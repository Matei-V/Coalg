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
	id_cap    int
	titlu     string
	descriere string
}

type pag struct {
	id_pag    int
	titlu     string
	descriere string
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

func Add_page(db *sql.DB, u unit, c chap, insert_chap bool, p pag, insert_pag bool, r []res, insert_res bool, pb []prob, insert_pb bool) {
	if insert_chap {
		insert, err := db.Query("insert into capitol(titlu, id_unit) values(?, ?)", c.titlu, u.id_unit)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	if insert_pag {
		id_c := strconv.Itoa(c.id_cap)
		insert, err := db.Query("insert into pagina(titlu, id_cap, descriere) values(?, ?, ?)", p.titlu, id_c, p.descriere)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	if insert_res {
		for i := 0; i < len(r); i++ {
			insert, err := db.Query("insert into resursa(titlu, link, limba, id_pag) values(?, ?, ?, ?)",
				r[i].titlu, r[i].link, r[i].limba, p.id_pag)
			if err != nil {
				panic(err.Error())
			}
			defer insert.Close()
		}
	}
	if insert_pb {
		for i := 0; i < len(pb); i++ {
			insert, err := db.Query("insert into problema(titlu, link, limba,  id_pag) values(?, ?, ?, ?)",
				pb[i].titlu, pb[i].link, pb[i].limba, p.id_pag)
			if err != nil {
				panic(err.Error())
			}
			defer insert.Close()
		}
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

func Get_Lessons(db *sql.DB, id_cap int) []string {
	var clmn []string
	row_num := Last_id(db, "pagina", "id_pag")
	for i := 0; i < row_num; i++ {
		var ans string
		stmt := "select titlu from pagina where id_pag=" + strconv.Itoa(i+1) + " and id_cap=" + strconv.Itoa(id_cap) + ";"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		if ans != "" {
			clmn = append(clmn, ans)
		}
	}
	return clmn
}

func Create_chap(id_cap int, titlu string) chap {
	var ans chap
	ans.id_cap = id_cap
	ans.titlu = titlu
	return ans
}

func Create_unit(id_unit int) unit {
	var ans unit
	ans.id_unit = id_unit
	return ans
}
func Create_pag(id_pag int, titlu string, descriere string) pag {
	var ans pag
	ans.id_pag = id_pag
	ans.titlu = titlu
	ans.descriere = descriere
	return ans
}

func Create_res(id_res []int, titlu []string /*, link []string, limba []string*/) []res {
	var ans []res

	for i := 0; i < len(id_res); i++ {
		var aux res
		aux.id_res = id_res[i]
		aux.titlu = titlu[i]
		aux.link = ""
		aux.limba = ""

		ans = append(ans, aux)
	}

	return ans
}

func Create_pb(id_pb []int, titlu []string /*, link []string, limba []string*/) []prob {
	var ans []prob

	for i := 0; i < len(id_pb); i++ {
		var aux prob
		aux.id_pb = id_pb[i]
		aux.titlu = titlu[i]
		aux.link = ""
		aux.limba = ""

		ans = append(ans, aux)
	}

	return ans
}

func Get_pag(db *sql.DB, titlu string) string {
	var ans string

	stmt := "select descriere from pagina where titlu='" + titlu + "';"
	row := db.QueryRow(stmt)
	row.Scan(&ans)
	return ans
}
