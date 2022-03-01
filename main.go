package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	types "github.com/strick-j/scimplistic/types"
	utils "github.com/strick-j/scimplistic/utils"
	views "github.com/strick-j/scimplistic/views"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}

func main() {

	logger := log.WithFields(log.Fields{
		"Category": "Server Processes",
		"Function": "main",
	})

	logger.Trace("Attempting to read configuration data")
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		logger.Error(err)
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
	}

	logger.Trace("Initializing mux router")
	r := mux.NewRouter()

	// Serve files for use, omit static from URL
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public/static/"))))

	// Handler for initial Index
	r.HandleFunc("/", views.IndexReq)

	// Handler for Settings functions
	setr := r.PathPrefix("/settings").Subrouter()
	setr.HandleFunc("/", views.SettingsHandler)
	setr.HandleFunc("/general", views.GeneralSettingsHandler)
	setr.HandleFunc("/secrets", views.SecretSettingsHandler)

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
	sr.HandleFunc("/{id}", views.SafeHandler).Methods("GET")
	sr.HandleFunc("/add", views.SafeAddHandler).Methods("POST")
	sr.HandleFunc("/update/{id}", views.SafeUpdateHandler).Methods("POST")
	sr.HandleFunc("/del/{id}", views.SafeDelHandler).Methods("POST")

	// Handlers for safe functions
	ar := r.PathPrefix("/accounts").Subrouter()
	ar.HandleFunc("/", views.AccountsHandler).Methods("GET")
	ar.HandleFunc("/{id}", views.AccountHandler).Methods("GET")
	ar.HandleFunc("/add", views.AccountAddHandler).Methods("POST")
	ar.HandleFunc("/update/{id}", views.AccountUpdateHandler).Methods("POST")
	ar.HandleFunc("/del/{id}", views.AccountDelHandler).Methods("POST")

	logger.Trace("mux router initialized")

	siteSettings.Router = r

	logger.Info("Attempting to start SCIMPLISTIC Server")

	utils.Start(&siteSettings)
}
