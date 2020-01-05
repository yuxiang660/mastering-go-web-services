package main

import (
    "fmt"
    "log"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func loopback(conn *websocket.Conn) {
	messageType, from, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

    for {	
		back := []byte(time.Now().String() + " : ")
		back = append(back, from...)
        if err := conn.WriteMessage(messageType, back); err != nil {
            log.Println(err)
            return
		}
		time.Sleep(time.Second)
    }
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// This will determine whether or not an incoming request from a different domain is allowed to
	// connect, and if it isn’t they’ll be hit with a CORS error.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	
	ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
		log.Println(err)
		return
	}

    loopback(ws)
}

func main() {
	fmt.Println("Hello World")
	
    http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/wsclient", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "hello-client.html")
	})

    http.ListenAndServe(":8080", nil)
}
