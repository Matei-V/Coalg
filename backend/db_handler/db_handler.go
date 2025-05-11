package db_handler

import (
	"database/sql"
	"fmt"
	"mysql_drive/utils"
	"strconv"
)

func Connect_db() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/coalgdb")
	if err != nil {
		fmt.Println("error")
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}

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

func Create_res(id_res []int, titlu []string, link []string) []res {
	var ans []res

	for i := 0; i < len(id_res); i++ {
		var aux res
		aux.id_res = id_res[i]
		aux.titlu = titlu[i]
		aux.link = link[i]

		ans = append(ans, aux)
	}

	return ans
}

func Create_pb(id_pb []int, titlu []string, link []string) []prob {
	var ans []prob

	for i := 0; i < len(id_pb); i++ {
		var aux prob
		aux.id_pb = id_pb[i]
		aux.titlu = titlu[i]
		aux.link = link[i]

		ans = append(ans, aux)
	}

	return ans
}

func Add_page(db *sql.DB, u unit, c chap, insert_chap bool, p pag, insert_pag bool, r []res, insert_res bool, pb []prob, insert_pb bool, autor string) {
	if insert_chap {
		insert, err := db.Query("insert into capitol(titlu, id_unit) values(?, ?)", c.titlu, u.id_unit)
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	if insert_pag {
		id_c := strconv.Itoa(c.id_cap)
		p.descriere = utils.Parse_code(p.descriere)
		insert, err := db.Query("insert into pagina(titlu, id_cap, descriere, autor, aprobat) values(?, ?, ?, ?, ?)", p.titlu, id_c, p.descriere, autor, "0")
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}
	if insert_res {
		for i := 0; i < len(r); i++ {
			fmt.Print("pag")
			fmt.Print(p.id_pag)
			insert, err := db.Query("insert into resursa(titlu, link, id_pag) values(?, ?, ?)",
				r[i].titlu, r[i].link, p.id_pag)
			if err != nil {
				panic(err.Error())
			}
			defer insert.Close()
		}
	}
	if insert_pb {
		for i := 0; i < len(pb); i++ {
			insert, err := db.Query("insert into problema(titlu, link,  id_pag) values(?, ?, ?)",
				pb[i].titlu, pb[i].link, p.id_pag)
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

func Get_Chapters(db *sql.DB, unit string) ([]string, []string) {
	var clmn, id []string
	row_num := Last_id(db, "capitol", "id_cap")

	for i := 0; i < row_num; i++ {
		var ans string
		stmt := "select titlu from capitol where id_cap=" + strconv.Itoa(i+1) + " and id_unit=" + unit + ";"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		if ans != "" {
			clmn = append(clmn, ans)
			id = append(id, strconv.Itoa(i+1))
		}
	}
	return clmn, id
}

func Get_Lessons(db *sql.DB, id_cap string) ([]string, []string) {
	var clmn, id []string
	row_num := Last_id(db, "pagina", "id_pag")
	for i := 0; i < row_num; i++ {
		var ans string
		stmt := "select titlu from pagina where id_pag=" + strconv.Itoa(i+1) + " and id_cap=" + id_cap + " and aprobat=1;"
		row := db.QueryRow(stmt)
		row.Scan(&ans)
		if ans != "" {
			clmn = append(clmn, ans)
			id = append(id, strconv.Itoa(i+1))
		}
	}
	return clmn, id
}

func Get_res(db *sql.DB, id_pag string) ([]string, []string) {
	var nr_res string
	var titlu_res, link []string
	var aux string
	stmt := "select count(*) from resursa where id_pag='" + id_pag + "';"
	row := db.QueryRow(stmt)
	row.Scan(&nr_res)

	n, _ := strconv.Atoi(nr_res)

	for i := 0; i < n; i++ {
		stmt = "select titlu from resursa where id_pag='" + id_pag + "' limit 1 offset " + strconv.Itoa(i) + ";"
		row = db.QueryRow(stmt)
		row.Scan(&aux)
		titlu_res = append(titlu_res, aux)

		stmt = "select link from resursa where id_pag='" + id_pag + "' limit 1 offset " + strconv.Itoa(i) + ";"
		row = db.QueryRow(stmt)
		row.Scan(&aux)
		link = append(link, aux)
	}

	return titlu_res, link

}

func Get_pb(db *sql.DB, id_pag string) ([]string, []string) {
	var nr_pb string
	var titlu_pb, link []string
	var aux string
	stmt := "select count(*) from problema where id_pag='" + id_pag + "';"
	row := db.QueryRow(stmt)
	row.Scan(&nr_pb)

	n, _ := strconv.Atoi(nr_pb)

	for i := 0; i < n; i++ {
		stmt = "select titlu from problema where id_pag='" + id_pag + "' limit 1 offset " + strconv.Itoa(i) + ";"
		row = db.QueryRow(stmt)
		row.Scan(&aux)
		titlu_pb = append(titlu_pb, aux)

		stmt = "select link from problema where id_pag='" + id_pag + "' limit 1 offset " + strconv.Itoa(i) + ";"
		row = db.QueryRow(stmt)
		row.Scan(&aux)
		link = append(link, aux)
	}

	return titlu_pb, link

}

func Get_title(db *sql.DB, id_pag string) string {
	var aux string
	stmt := "select titlu from pagina where id_pag='" + id_pag + "';"
	row := db.QueryRow(stmt)
	row.Scan(&aux)
	return aux
}

func Get_desc(db *sql.DB, id_pag string) string {
	var ans string

	stmt := "select descriere from pagina where id_pag='" + id_pag + "';"
	row := db.QueryRow(stmt)
	row.Scan(&ans)
	return ans
}
