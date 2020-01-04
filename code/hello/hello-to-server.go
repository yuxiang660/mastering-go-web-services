package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("word"))
}

func main() {

	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/send", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.FormValue("word"))
	})

	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":8080", nil)
}
