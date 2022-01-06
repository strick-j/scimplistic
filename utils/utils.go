package utils

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

var httpAddr string = ":8080"
var httpsAddr string = ":8443"

func GetRedirectUrl(referer string) string {
	var redirectUrl string
	url := strings.Split(referer, "/")

	if len(url) > 4 {
		redirectUrl = "/" + strings.Join(url[3:], "/")
	} else {
		redirectUrl = "/"
	}
	return redirectUrl
}

func startTlsListen(certPath string, keyPath string) {
	log.Printf("INFO startTlsListen: Attempting to start TLS Server.")
	srv := http.Server{
		Addr: httpsAddr,
	}

	_, tlsPort, err := net.SplitHostPort(httpsAddr)
	if err != nil {
		return
	}
	go redirectToHTTPS(tlsPort)

	srv.ListenAndServeTLS(certPath, keyPath)
}

func redirectToHTTPS(tlsPort string) {
	log.Print("INFO RedirectToHTTPS: Attempting to create redirect")
	httpSrv := http.Server{
		Addr: httpAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, _ := net.SplitHostPort(r.Host)
			u := r.URL
			u.Host = net.JoinHostPort(host, tlsPort)
			u.Scheme = "https"
			log.Println(u.String())
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		}),
	}
	log.Println("INFO RedirectToHTTPS: Listening on:", httpSrv.ListenAndServe())
}

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}
