package main

import (
	"fmt"
	"net/http"
	"log"
	"sync"
)

const (
	serverName = "localhost"
	SSLport = ":8181"
	HTTPport = ":8080"
	SSLprotocol = "https://"
	HTTPprotocol = "http://"
)

func secureRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have arrived at port 443, and now you are marginally more secure.")
}

func redirectNonSecure(w http.ResponseWriter, r *http.Request) {
	log.Println("Non-secure request initiated, redirecting.")
	redirectURL := SSLprotocol + serverName + SSLport + r.RequestURI
	http.Redirect(w, r, redirectURL, http.StatusOK)
}

func main() {
	wg := sync.WaitGroup{}
	log.Println("Starting redirection server, try to access @ http:")

	wg.Add(1)
	go func(){
		http.ListenAndServe(HTTPport, http.HandlerFunc(redirectNonSecure))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		http.ListenAndServeTLS(SSLport, "cert.pem", "key.pem", http.HandlerFunc(secureRequest))
		wg.Done()
	}()

	wg.Wait()
}
