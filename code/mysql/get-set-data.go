package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

var database *sql.DB

type User struct {
	ID int "json:id"
	Name string "json:username"
	Email string "json:email"
	First string "json:first"
	Last string "json:last"
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")
	
	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Something went wrong!")
	}

	sql := "INSERT INTO users set user_nickname='" + NewUser.Name +
		   "', user_first='" + NewUser.First +
		   "', user_last='" + NewUser.Last +
		   "', user_email='" + NewUser.Email + "'"
	
	q, err := database.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(q)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	id := urlParams["id"]
	ReadUser := User{}
	err := database.QueryRow("select * from users where user_id=?", id).Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.First, &ReadUser.Last, &ReadUser.Email)
	
	switch{
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No such user")
	case err != nil:
		fmt.Fprintf(w, "Error")
		log.Fatal(err)
	default:
		output, _ := json.Marshal(ReadUser)
		fmt.Fprintf(w, string(output))
	}
}

func main() {
	db, err := sql.Open("mysql", "ben:123456@/social_network")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/api/user/create", CreateUser)
	routes.HandleFunc("/api/user/read/{id:[0-9]+}", GetUser)

	http.Handle("/", routes)
	http.ListenAndServe(":8080", nil)
}