package api

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"mysql_drive/api/auth"
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
	caps, id_cap := dbh.Get_Chapters(db, unit)
	var htmlStr string
	currentId := param.Get("id")

	if len(caps) == 0 {
		htmlStr += "<div class='menu-empty'><p>Nu există încă lecții pentru această clasă.</p><a href='/add-page' class='btn-add-mini'>+ Adaugă prima lecție</a></div>"
	}

	for i := 0; i < len(caps); i++ {
		if caps[i] == "" {
			continue
		}
		less, id := dbh.Get_Lessons(db, id_cap[i])
		htmlStr += "<div class='menu-chapter-group'>"
		htmlStr += "<div class='menu-chapter-header'>"
		htmlStr += "<span class='chapter-icon'>📂</span>"
		htmlStr += "<h4 class='capitol'>" + template.HTMLEscapeString(caps[i]) + "</h4>"
		htmlStr += "<span class='lesson-count-badge'>" + strconv.Itoa(len(less)) + "</span>"
		htmlStr += "</div>"

		htmlStr += "<div class='menu-lesson-list'>"
		for j := 0; j < len(less); j++ {
			paramValues := url.Values{}
			paramValues.Add("clasa", unit)
			paramValues.Add("id", id[j])
			activeClass := ""
			if currentId == id[j] {
				activeClass = " active-lesson"
			}
			htmlStr += "<div class='less" + activeClass + "'><a href='/unit?" + paramValues.Encode() + "'>"
			htmlStr += "<span class='less-bullet'>•</span>"
			htmlStr += "<span class='less-text'>" + template.HTMLEscapeString(less[j]) + "</span>"
			htmlStr += "</a></div>"
		}
		if len(less) == 0 {
			htmlStr += "<div class='less-empty'>Nicio lecție publicată încă</div>"
		}
		htmlStr += "</div></div>"
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
	id := param.Get("id")
	unit := param.Get("clasa")
	if id == "" {
		htmlStr := "<div id='content' class='lesson-empty-container'>"
		htmlStr += "<div class='lesson-empty-card'>"
		htmlStr += "<div class='empty-card-icon'>📖</div>"
		htmlStr += "<h3>Alege o temă pentru a începe studiul</h3>"
		if unit != "" {
			htmlStr += "<p>Navighează prin capitolele pentru <strong>Clasa a " + template.HTMLEscapeString(unit) + "-a</strong> din panoul lateral și alege o lecție pentru a citi explicațiile și a rezolva problemele propuse.</p>"
		} else {
			htmlStr += "<p>Alege o clasă și o temă din meniu pentru a citi explicațiile teoretice, codul sursă C++ și problemele recomandate.</p>"
		}
		htmlStr += "<div class='empty-card-actions'>"
		htmlStr += "<a href='/add-page' class='btn btn-primary-gradient'>+ Propune o lecție nouă</a>"
		htmlStr += "</div>"
		htmlStr += "</div>"
		htmlStr += "</div>"

		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
		return
	}

	desc := dbh.Get_desc(db, id)
	title := dbh.Get_title(db, id)
	titlu_res, link_res := dbh.Get_res(db, id)
	titlu_pb, link_pb := dbh.Get_pb(db, id)

	htmlStr := "<div id='content' class='lesson-container'>"
	htmlStr += "<div class='lesson-header'>"
	htmlStr += "<div class='lesson-breadcrumbs'>"
	if unit != "" {
		htmlStr += "<a href='/'>Acasă</a> <span class='sep'>/</span> <a href='/unit?clasa=" + template.HTMLEscapeString(unit) + "'>Clasa a " + template.HTMLEscapeString(unit) + "-a</a> <span class='sep'>/</span> <span class='current'>" + template.HTMLEscapeString(title) + "</span>"
	} else {
		htmlStr += "<a href='/'>Acasă</a> <span class='sep'>/</span> <span class='current'>" + template.HTMLEscapeString(title) + "</span>"
	}
	htmlStr += "</div>"
	htmlStr += "<div class='lesson-title-row'>"
	htmlStr += "<h1 class='lesson-title'>" + template.HTMLEscapeString(title) + "</h1>"
	if unit != "" {
		htmlStr += "<span class='badge badge-grade'>Clasa a " + template.HTMLEscapeString(unit) + "-a</span>"
	}
	htmlStr += "</div>"
	htmlStr += "</div>"

	htmlStr += "<div class='lesson-body-card'>"
	htmlStr += "<div class='lesson-prose'>" + desc + "</div>"
	htmlStr += "</div>"

	if len(titlu_res) > 0 {
		htmlStr += "<div class='lesson-section'>"
		htmlStr += "<h3 class='section-title'><span class='section-icon'>📚</span> Resurse Utile & Documentație</h3>"
		htmlStr += "<div class='resource-grid'>"
		for i := 0; i < len(titlu_res); i++ {
			htmlStr += "<a class='resource-card' href='" + template.HTMLEscapeString(link_res[i]) + "' target='_blank' rel='noopener noreferrer'>"
			htmlStr += "<div class='resource-icon'>🔗</div>"
			htmlStr += "<div class='resource-info'>"
			htmlStr += "<div class='resource-name'>" + template.HTMLEscapeString(titlu_res[i]) + "</div>"
			htmlStr += "<div class='resource-url'>" + template.HTMLEscapeString(link_res[i]) + "</div>"
			htmlStr += "</div>"
			htmlStr += "<span class='external-arrow'>↗</span>"
			htmlStr += "</a>"
		}
		htmlStr += "</div></div>"
	}

	if len(titlu_pb) > 0 {
		htmlStr += "<div class='lesson-section'>"
		htmlStr += "<h3 class='section-title'><span class='section-icon'>🎯</span> Probleme Recomandate</h3>"
		htmlStr += "<div class='problem-grid'>"
		for i := 0; i < len(titlu_pb); i++ {
			htmlStr += "<a class='problem-card' href='" + template.HTMLEscapeString(link_pb[i]) + "' target='_blank' rel='noopener noreferrer'>"
			htmlStr += "<div class='problem-badge'>Pbinfo / Concurs</div>"
			htmlStr += "<div class='problem-name'>" + template.HTMLEscapeString(titlu_pb[i]) + "</div>"
			htmlStr += "<div class='problem-action'>Rezolvă problema <span>→</span></div>"
			htmlStr += "</a>"
		}
		htmlStr += "</div></div>"
	}

	htmlStr += "</div>"
	tmpl, _ := template.New("t").Parse(htmlStr)
	tmpl.Execute(w, nil)
}

func Auth_cont(w http.ResponseWriter, r *http.Request) {
	if err := auth.Authorize(r); err == nil {
		htmlStr := "static/template/acc_auth_cont.html"
		tmpl, _ := template.ParseFiles(htmlStr)
		tmpl.Execute(w, nil)
		return
	}
	htmlStr := "static/template/acc_unauth_cont.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, nil)
}

func MyAcc(w http.ResponseWriter, r *http.Request) {
	username, _ := r.Cookie("username")
	var less []string
	var id_unit, id_less, aprobat []string
	if dbh.IsAdmin(db, username.Value) == false {
		less, id_unit, id_less, aprobat = dbh.Get_Lessons_Auth(db, username.Value)
	} else {
		less, id_unit, id_less, aprobat = dbh.Get_Lessons_Admin(db, username.Value)
	}
	type Input struct {
		Less    string
		Cap     string
		ID_Less string
		Aprobat string
	}
	var In []Input
	fmt.Print("less:" + strconv.Itoa(len(less)) + "\n")
	fmt.Print("id_unit:" + strconv.Itoa(len(id_unit)) + "\n")
	fmt.Print("id_less:" + strconv.Itoa(len(id_less)) + "\n")
	fmt.Print("aprobat:" + strconv.Itoa(len(aprobat)) + "\n")
	for i := 0; i < len(less); i++ {
		In = append(In, Input{less[i], id_unit[i], id_less[i], aprobat[i]})
	}
	htmlStr := "static/template/less_list.html"
	tmpl, _ := template.ParseFiles(htmlStr)
	tmpl.Execute(w, In)
}

func Aproba(w http.ResponseWriter, r *http.Request) {
	my_url, _ := url.Parse(r.URL.String())
	param, _ := url.ParseQuery(my_url.RawQuery)
	dbh.Aprove(db, param.Get("id"))
}
func Respinge(w http.ResponseWriter, r *http.Request) {
	my_url, _ := url.Parse(r.URL.String())
	param, _ := url.ParseQuery(my_url.RawQuery)
	dbh.Reject(db, param.Get("id"))
}
