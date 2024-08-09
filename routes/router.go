package routes

import (
	"amanks/voting/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("POST")

	// for websocket
	r.HandleFunc("/ws", controllers.HandleConnection)
	r.HandleFunc("/vote", controllers.VoteHandler).Methods("POST")
	r.HandleFunc("/results", controllers.ResultsHandler).Methods("GET")

	return r
}
