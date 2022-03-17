package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"video-chat-app/server"

	"github.com/joho/godotenv"
)

func main() {
	server.AllRooms.Init()
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error in loading env file %s", err.Error())
	}

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	port := os.Getenv("PORT")
	if !(port == "") {
		log.Println("Starting server on port" + port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Printf("Error in starting server %s", err.Error())
		}
	} else {
		log.Println("Starting server on port 8000")
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			fmt.Printf("Error in starting server %s", err.Error())
		}
	}
}
