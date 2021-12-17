package main

import (
	"net/http"

	"github.com/strick-j/go-form-webserver/views"
)

func main() {

	// Handler for initial Index
	http.HandleFunc("/", views.IndexReq)

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
	http.HandleFunc("/safedelform/", views.SafeDelForm)

	// Serve files for use, omit static from URL
	http.Handle("/static/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":8080", nil)
}
