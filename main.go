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
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))

	// Handler for initial Index
	r.HandleFunc("/", views.IndexReq)

	// Handler for Settings functions
	r.HandleFunc("/settings/", views.SettingsHandler)
	r.HandleFunc("/settings/{type}", views.SettingsTypeHandler)

	// Handlers for user functions
	r.HandleFunc("/users/", views.UsersHandler)
	r.HandleFunc("/users/{action}", views.UsersActionHandler)

	// Handlers for group functions
	r.HandleFunc("/groups/", views.GroupsHandler)
	r.HandleFunc("/groups/{action}", views.GroupsActionHandler)
	r.HandleFunc("/groups/{action}/{id}", views.GroupsActionHandler)

	// Handlers for safe functions
	r.HandleFunc("/safes/", views.SafesHandler)
	r.HandleFunc("/safes/{action}", views.SafesActionHandler)

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
