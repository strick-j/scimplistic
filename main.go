package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/strick-j/scimplistic/config"
	"github.com/strick-j/scimplistic/views"
)

var httpAddr string = ":8080"
var httpsAddr string = ":8443"

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
	http.HandleFunc("/userdel/", views.UserDelFunc)

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
	http.HandleFunc("/safedel/", views.SafeDelFunc)

	// Read in config values
	values, err := config.ReadConfig("config.json")
	if err != nil {
		log.Println("ERROR IndexReq:", err)
	} else if values.EnableHTTPS {
		log.Printf("INFO Main: HTTPS enabled, checking cert and key files necessary to serve HTTPS")
		if values.CertName != "" {
			log.Println("INFO Main: Certificate Name Identified:", values.CertName)
		} else {
			log.Println("ERROR Main: Certificate Name not Found. Serving standard HTTP")
			return
		}

		if values.KeyName != "" {
			log.Println("INFO Main: Key Name Identified:", values.KeyName)
		} else {
			log.Println("ERROR Main: Key Name not Found. Serving standard HTTP")
			return
		}

		srv := http.Server{
			Addr: httpsAddr,
		}

		_, tlsPort, err := net.SplitHostPort(httpsAddr)
		if err != nil {
			return
		}
		go utils.redirectToHTTPS(tlsPort)

		srv.ListenAndServeTLS(values.CertName, values.KeyName)

	} else {
		log.Printf("INFO Main: HTTPS not configured, serving standard HTTP")
		httpServerExitDone := &sync.WaitGroup{}

		utils.startHttpServer(httpServerExitDone)
	}

}
