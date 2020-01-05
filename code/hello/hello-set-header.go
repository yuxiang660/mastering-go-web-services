package main

import (
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Pragma", "no-cache")
		fmt.Fprintf(w, "Hello")
	})

	http.ListenAndServe(":8080", nil)
}