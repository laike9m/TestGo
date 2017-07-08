package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var logger *log.Logger

type Person struct {
	Name string
	Age  int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	indexFile, err := os.Open("index.html")
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			if string(msg) == "ping" {
				fmt.Println("ping")
				time.Sleep(2 * time.Second)
				err = conn.WriteMessage(msgType, []byte("pong"))
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				conn.Close()
				fmt.Println(string(msg))
				return
			}
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, string(index)) })

	// http.ListenAndServe(":3000", nil)
	err = http.ListenAndServeTLS(
		":3000",
		"cert.pem",
		"private_key",
		nil,
	)
	log.Fatal(err)
}
