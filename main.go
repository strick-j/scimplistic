package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	config "github.com/strick-j/scimplistic/config"
	types "github.com/strick-j/scimplistic/types"
	utils "github.com/strick-j/scimplistic/utils"
	views "github.com/strick-j/scimplistic/views"
)

func main() {

	r := mux.NewRouter()

	// Serve files for use, omit static from URL
	//r.Handle("/static/{rest}", http.StripPrefix("/static/", http.FileServer(http.Dir("public/static/"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))

	// Handler for initial Index
	r.HandleFunc("/", views.IndexReq)

	// Handler for Settings functions
	r.HandleFunc("/settings/", views.SettingsForm)
	r.HandleFunc("/configuresettings/", views.ConfigureSettings)

	// Handlers for user functions
	r.HandleFunc("/allusers/", views.UserAllReq)
	r.HandleFunc("/useraddform/", views.UserAddForm)
	r.HandleFunc("/useraddreq/", views.UserAddReq)
	r.HandleFunc("/userdel/", views.UserDelFunc)

	// Handlers for group functions
	r.HandleFunc("/allgroups/", views.GroupAllReq)
	r.HandleFunc("/groupaddform/", views.GroupAddForm)
	r.HandleFunc("/groupaddreq/", views.GroupAddReq)
	r.HandleFunc("/groupdel/", views.GroupDelFunc)
	r.HandleFunc("/groupupdate/", views.GroupUpdateForm)
	r.HandleFunc("/groupupdatereq/", views.GroupUpdateFunc)

	// Handlers for safe functions
	r.HandleFunc("/allsafes/", views.SafeAllReq)
	r.HandleFunc("/safeaddform/", views.SafeAddForm)
	r.HandleFunc("/safeaddreq/", views.SafeAddReq)
	r.HandleFunc("/safedel/", views.SafeDelFunc)

	values, err := config.ReadConfig("settings.json")
	if err != nil {
		log.Println("ERROR Main:", err)
	}

	siteSettings := types.ConfigSettings{
		ServerName:     values.HostName,
		MaxConnections: values.MaxConnections,
		HostName:       values.HostName,
		HostAlias:      values.HostAlias,
		IP:             values.IP,
		Port:           values.Port,
		TLS:            values.TLS,
		CertFile:       values.CertFile,
		PrivKeyFile:    values.PrivKeyFile,
		Router:         r,
	}

	log.Printf("INFO MAIN: Attempting to start Scimplistic server")

	utils.Start(&siteSettings)

}
