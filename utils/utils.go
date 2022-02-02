package utils

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/strick-j/scimplistic/types"
)

var (
	httpServer *http.Server

	settings *types.ConfigSettings

	serverStarted  bool       = false
	serverPaused   bool       = false
	serverStopping bool       = false
	serverEndChan  chan error = make(chan error)

	startCallback  func()
	pauseCallback  func()
	stopCallback   func()
	resumeCallback func()

	//SERVER VERSION NUMBER
	version string = "1.0-BETA.0"
)

func Start(s *types.ConfigSettings) {
	log.Printf("INFO Start: Checking Server Status")
	if serverStarted || serverPaused {
		log.Printf("INFO Start: Server running or paused")
		return
	}
	serverStarted = true
	log.Println("INFO Start: Starting server...")
	// Set server settings
	if s != nil {
		settings = s
	} else {
		// Default localhost settings
		log.Println("INFO Start: Using default settings...")
		settings = &types.ConfigSettings{
			ServerName:     "!server!",
			MaxConnections: 0,
			HostName:       "localhost",
			HostAlias:      "localhost",
			IP:             "localhost",
			Port:           8080,
			TLS:            false,
			CertFile:       "",
			PrivKeyFile:    "",
			OriginOnly:     false}
	}

	// Start socket listener
	httpServer = makeServer("/", settings.TLS, s.Router)

	// Run callback
	if startCallback != nil {
		startCallback()
	}

	// Start macro listener
	go macroListener()

	log.Println("INFO Start: Startup complete")

	// Wait for server shutdown
	doneErr := <-serverEndChan

	if doneErr != http.ErrServerClosed {
		log.Println("ERROR Start: Fatal server error:", doneErr.Error())

	}

	log.Println("INFO Start: Server shut-down completed")

	serverStarted = false

	if stopCallback != nil {
		log.Println("executing stop callback")
		stopCallback()
	}

}

func makeServer(handleDir string, tls bool, router *mux.Router) *http.Server {
	server := &http.Server{
		Addr:         settings.IP + ":" + strconv.Itoa(settings.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	if tls {
		go func() {
			log.Println("INFO makeServer: Attempting to start HTTPS server on IP:" + settings.IP + " and Port:" + strconv.Itoa(settings.Port))
			err := server.ListenAndServeTLS(settings.CertFile, settings.PrivKeyFile)
			serverEndChan <- err
		}()
	} else {
		go func() {
			log.Println("INFO makeServer: Attempting to start server on IP:" + settings.IP + " and Port:" + strconv.Itoa(settings.Port))
			err := server.ListenAndServe()
			serverEndChan <- err
		}()
	}

	//
	return server
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
//   Server actions   ////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////

// Pause will temporarily pause the server
func Pause() {
	if !serverPaused {
		serverPaused = true

		SetPauseCallback(serverPaused)

		// Run callback
		if pauseCallback != nil {
			log.Println("running some callback", &pauseCallback)
			pauseCallback()
		}

		log.Println("Server paused")

		serverStarted = false
	}

}

// Resume will resume after pause
func Resume() {
	if serverPaused {
		serverStarted = true

		// Run callback
		if resumeCallback != nil {
			resumeCallback()
		}

		log.Println("Server resumed")

		serverPaused = false
	}
}

// ShutDown will shut the server down
func ShutDown() error {
	if !serverStopping {
		serverStopping = true

		// Shut server down
		log.Println("INFO Shutdown: Shutting server down...")
		shutdownErr := httpServer.Shutdown(context.Background())
		if shutdownErr != http.ErrServerClosed {
			return shutdownErr
		}
	}
	return nil
}

func ServerPaused() {
	log.Println("Establishing server Paused.")
}

func ServerStopped(action string) {
	log.Println("Establishing server Stopped.")
}
