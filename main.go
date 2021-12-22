package main

import (
	"net/http"

	"github.com/strick-j/go-form-webserver/views"
)

func main() {

	// Serve files for use, omit static from URL
	http.Handle("/static/", http.FileServer(http.Dir("public")))

	// Handler for initial Index
	http.HandleFunc("/", views.IndexReq)

	// Handler for Settings functions
	http.HandleFunc("/settings/", views.SettingsForm)
	http.HandleFunc("/configuresettings/", views.ConfigureSettings)

	// Handlers for user functions
	http.HandleFunc("/allusers/", views.UserAllReq)
	http.HandleFunc("/useraddform/", views.UserAddForm)
	http.HandleFunc("/useraddreq/", views.UserAddReq)
	http.HandleFunc("/userdelform/", views.UserDelForm)

	// Handlers for group functions
	http.HandleFunc("/allgroups/", views.GroupAllReq)
	http.HandleFunc("/groupaddform/", views.GroupAddForm)
	http.HandleFunc("/groupaddreq/", views.GroupAddReq)
	http.HandleFunc("/groupdel/", views.GroupDelFunc)
	http.HandleFunc("/groupupdate/", views.GroupUpdateForm)
	http.HandleFunc("/groupupdatereq/", views.GroupUpdateFunc)

	// Handlers for safe functions
	http.HandleFunc("/allsafes/", views.SafeAllReq)
	http.HandleFunc("/safeaddform/", views.SafeAddForm)
	http.HandleFunc("/safeaddreq/", views.SafeAddReq)
	http.HandleFunc("/safedelform/", views.SafeDelForm)

	http.ListenAndServe(":8080", nil)
}
