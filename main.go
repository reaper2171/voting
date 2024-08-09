package main

import (
	"amanks/voting/controllers"
	"amanks/voting/routes"
	"log"
	"net/http"
)

func main() {
	// setting up routes
	r := routes.SetupRoutes()

	// goroutine to handle votes in realtime(couting and saving their values in redis)
	go controllers.HandleVote()

	log.Print("Server starting at 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
