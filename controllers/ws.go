package controllers

import (
	"amanks/voting/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Vote)

// var redisPool = config.NewPool()
var upgrader = websocket.Upgrader{}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error handling ws conn: %v", err)
	}
	defer ws.Close()

	clients[ws] = true
	for {
		var vote models.Vote
		err := ws.ReadJSON(&vote)
		if err != nil {
			log.Printf("error reading vote: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- vote
	}
}

func HandleVote() {
	for {
		vote := <-broadcast
		conn := redisPool.Get()
		defer conn.Close()
		_, err := conn.Do("HINCRBY", vote.Session, vote.Opt, 1)
		if err != nil {
			log.Printf("error updating session vote value: %v", err)
		}
		for client := range clients {
			err := client.WriteJSON(vote)
			if err != nil {
				log.Printf("error sending votes in broadcast: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
