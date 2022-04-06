package utils

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/internal/types"
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
	version string = "1.0-BETA.1"
)

func Start(s *types.ConfigSettings) {
	logger := log.WithFields(log.Fields{
		"Category": "Server Processes",
		"Function": "Start",
	})

	logger.Info("Checking Server Status")
	if serverStarted || serverPaused {
		logger.Info("INFO Start: Server running or paused")
		return
	}
	serverStarted = true
	logger.Info("Starting server...")
	// Set server settings
	if s != nil {
		settings = s
	} else {
		// Default localhost settings
		logger.Info("Using default settings...")
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

	logger.Info("Startup complete")

	// Wait for server shutdown
	doneErr := <-serverEndChan

	if doneErr != http.ErrServerClosed {
		logger.Error("Fatal server error:", doneErr.Error())

	}

	logger.Info("Server shut-down completed")

	serverStarted = false

	if stopCallback != nil {
		logger.Info("executing stop callback")
		stopCallback()
	}

}

func makeServer(handleDir string, tls bool, router *mux.Router) *http.Server {
	logger := log.WithFields(log.Fields{
		"Category": "Server Processes",
		"Function": "makeServer",
	})
	server := &http.Server{
		Addr:         settings.IP + ":" + strconv.Itoa(settings.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	if tls {
		go func() {
			logger.Info("Attempting to start HTTPS server on IP:" + settings.IP + " and Port:" + strconv.Itoa(settings.Port))
			err := server.ListenAndServeTLS(settings.CertFile, settings.PrivKeyFile)
			serverEndChan <- err
		}()
	} else {
		go func() {
			logger.Info("Attempting to start server on IP:" + settings.IP + " and Port:" + strconv.Itoa(settings.Port))
			err := server.ListenAndServe()
			serverEndChan <- err
		}()
	}

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
			pauseCallback()
		}

		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "Pause"}).Info("Server paused")

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

		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "Resume"}).Info("Server resumed")

		serverPaused = false
	}
}

// ShutDown will shut the server down
func ShutDown() error {
	if !serverStopping {
		serverStopping = true

		// Shut server down
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "ShutDown"}).Info("Shutting server down...")
		shutdownErr := httpServer.Shutdown(context.Background())
		if shutdownErr != http.ErrServerClosed {
			return shutdownErr
		}
	}
	return nil
}

func ServerPaused() {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "ServerPaused"}).Info("Establishing server Paused.")
}

func ServerStopped(action string) {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "ServerStopped"}).Info("Establishing server Stopped.")
}

//////// Generic Utilities /////////////////////////////////////////////////////////////////////////////////////

// ReadConfig will read the configuration json file to read the parameters
// which will be passed in the config file
func ReadConfig(fileName string) (config types.ConfigSettings, err error) {

	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "ReadConfig"}).Error("Unable to read config file, switching to flag mode")
		return types.ConfigSettings{}, err
	}
	//log.Print(configFile)
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "ReadConfig"}).Error("Invalid JSON")
		return types.ConfigSettings{}, err
	}
	return config, nil
}
