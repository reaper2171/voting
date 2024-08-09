package controllers

import (
	"amanks/voting/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	var vote models.Vote
	json.NewDecoder(r.Body).Decode(&vote)
	log.Printf("vote value: %v", vote)
	broadcast <- vote
	w.WriteHeader(http.StatusOK)
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session")
	conn := redisPool.Get()
	defer conn.Close()

	results, err := redis.String(conn.Do("HGETALL", session))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
