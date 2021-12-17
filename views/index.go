package views

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func indexRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	/*res := BuildUrl("Users", "GET")

	var bodyObject types.User
	json.Unmarshal(res, &bodyObject)

	*/
	tpl.ExecuteTemplate(w, "index.html", nil)
}
