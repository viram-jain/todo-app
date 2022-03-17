package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// AllRooms global hashmap for the server
var AllRooms RoomMap

// CreateRoomRequestHandler Create a room and return room ID
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	err := json.NewEncoder(w).Encode(resp{RoomID: roomID})
	if err != nil {
		fmt.Printf("CreateRoomRequestHandler error %s", err.Error())
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast

		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					fmt.Printf("broadcaster error WriteJSON error %s", err.Error())
					client.Conn.Close()
					break
				}
			}
		}
	}
}

// JoinRoomRequestHandler Join in a particular room
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		fmt.Println("JoinRoomRequestHandler error : roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("JoinRoomRequestHandler error : Upgrade error:", err.Error())
		return
	}
	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			fmt.Println("JoinRoomRequestHandler error ReadJSON error", err.Error())
			break
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		broadcast <- msg
	}

}
