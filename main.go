package main

import (
	"html/template"
	"net/http"

	"github.com/strick-j/go-form-webserver/views"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {

	http.HandleFunc("/", index)

	// Handlers for user functions
	http.HandleFunc("/allusers/", views.UserAllReq)
	http.HandleFunc("/useraddform/", views.UserAddForm)
	http.HandleFunc("/useraddreq/", views.UserAddReq)
	http.HandleFunc("/userdelform/", views.UserDelForm)

	// Handlers for group functions
	http.HandleFunc("/allgroups/", views.GroupAllReq)
	http.HandleFunc("/groupaddform/", views.GroupAddForm)
	http.HandleFunc("/groupaddreq/", views.GroupAddReq)
	http.HandleFunc("/groupdelform/", views.GroupDelForm)

	// Handlers for safe functions
	http.HandleFunc("/allsafes/", views.SafeAllReq)
	http.HandleFunc("/safeaddform/", views.SafeAddForm)
	http.HandleFunc("/safeaddreq/", views.SafeAddReq)
	//http.HandleFunc("/safedelform/", views.SafeDelForm)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}
