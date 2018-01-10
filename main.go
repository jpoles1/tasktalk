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

func init() {
	gotenv.Load()
	mongoLoad()
}
func main() {
	//Request routing
	router := mux.NewRouter()
	router.HandleFunc("/hello", hello).Methods("GET")
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
