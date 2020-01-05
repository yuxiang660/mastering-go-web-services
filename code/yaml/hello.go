package main

import (
	"net/http"
	"fmt"
	"gopkg.in/yaml.v2"
)

type User struct {
	Name string `yaml:"name"`
	Email string `yaml:"email"`
	ID int `yaml:"id"`
}

func userRouter(w http.ResponseWriter, r *http.Request) {
	ourUser := User{}
	ourUser.Name = "Bill Smith"
	ourUser.Email = "bill@example.com"
	ourUser.ID = 100

	output, _ := yaml.Marshal(&ourUser)
	fmt.Fprintln(w, string(output))
}

func main() {
	fmt.Println("Starting YAML server")
	http.HandleFunc("/user", userRouter)
	http.ListenAndServe(":8080", nil)
}