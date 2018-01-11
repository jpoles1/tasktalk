package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var fbToken string

func init() {
	gotenv.Load()
	fbToken = os.Getenv("PAGE_ACCESS_TOKEN")
	if fbToken == "" {
		log.Fatal("No FB Access Token supplied in .env config file!")
	}
	dbLoad()
}
func main() {
	//Request routing
	router := mux.NewRouter()
	router.HandleFunc("/", hello).Methods("GET").Queries("hub.challenge", "{hub.challenge}")
	router.HandleFunc("/", receiveMsg).Methods("POST")
	//Start the engines
	portPtr := flag.String("p", "3333", "Server Port")
	flag.Parse()
	port := ":" + *portPtr
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}
	color.Green("Starting server on port: %s", port[1:])
	color.Green("Access server locally at: http://127.0.0.1:%s", port[1:])
	//Handling system signals
	log.Fatal(http.ListenAndServe(port, router))
	fmt.Println("Terminating TaskTalk Server...")
}
