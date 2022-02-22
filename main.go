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
	setr := r.PathPrefix("/settings").Subrouter()
	setr.HandleFunc("/", views.SettingsHandler)
	setr.HandleFunc("/general", views.SettingsGenHandler)
	setr.HandleFunc("/secrets", views.SettingsSecretHandler)

	// Handlers for user functions
	ur := r.PathPrefix("/users").Subrouter()
	ur.HandleFunc("/", views.UsersHandler).Methods("GET")
	ur.HandleFunc("/{id}", views.UserHandler).Methods("GET")
	ur.HandleFunc("/add", views.UserAddHandler).Methods("POST")
	ur.HandleFunc("/update/{id}", views.UserUpdateHandler).Methods("POST")
	ur.HandleFunc("/del/{id}", views.UserDelHandler).Methods("POST")

	// Handlers for group functions
	gr := r.PathPrefix("/groups").Subrouter()
	gr.HandleFunc("/", views.GroupsHandler).Methods("GET")
	gr.HandleFunc("/{id}", views.GroupHandler).Methods("GET")
	gr.HandleFunc("/add", views.GroupAddHandler).Methods("POST")
	gr.HandleFunc("/update/{id}", views.GroupUpdateHandler).Methods("POST")
	gr.HandleFunc("/del/{id}", views.GroupDelHandler).Methods("POST")

	// Handlers for safe functions
	sr := r.PathPrefix("/safes").Subrouter()
	sr.HandleFunc("/", views.SafesHandler).Methods("GET")
	gr.HandleFunc("/{id}", views.SafeHandler).Methods("GET")
	sr.HandleFunc("/add", views.SafeAddHandler).Methods("POST")
	sr.HandleFunc("/update/{id}", views.SafeUpdateHandler).Methods("POST")
	sr.HandleFunc("/del/{id}", views.SafeDelHandler).Methods("POST")

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
