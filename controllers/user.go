package controllers

import (
	"amanks/voting/config"
	"amanks/voting/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

// var loginPageURL = "http://localhost:8080/loginpage"
var loginPageURL = "https://www.google.com/"

// var HomePageURl = "https://locahost:8080/homepage"
var HomePageURl = "https://www.google.com/"

var redisPool = config.NewPool()

// API for registering a new user and saving its value in redis
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Bad Request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}
	conn := redisPool.Get()
	_, err = conn.Do("HSET", "user-map", req.Username, req.Password)
	if err != nil {
		log.Printf("error saving  username and pass in redis: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	log.Printf("User register with username %s and password %s", req.Username, req.Password)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Registerd"))
}

// API for logging in with username and password and generating JWT tokens
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.RegReq
	json.NewDecoder(r.Body).Decode(&user)
	log.Printf("The user request value %v", user.Username)

	conn := redisPool.Get()
	userPass, err := redis.String(conn.Do("HGET", "user-map", user.Username))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("No username found: %v", err)
		w.Write([]byte("No such user found"))
		return
	}

	if user.Password == userPass {
		tokenString, err := config.CreateToken(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print("No username found")
			return
		}
		_, err = conn.Do("HSET", "token-map", user.Username, tokenString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error saving token in redis-map: %v", err)
			return
		}
		log.Print("control flow is heres")
		http.Redirect(w, r, HomePageURl, http.StatusPermanentRedirect)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print(w, "Invalid credentials")
		return
	}
}

// API for logging out
func Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Printf("error decoding logout req: %v", err)
		return
	}
	conn := redisPool.Get()
	_, err = conn.Do("DEL", "token-map", req.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error deleting token from redis: %v", err)
		return
	}

	http.Redirect(w, r, loginPageURL, 200)
}
