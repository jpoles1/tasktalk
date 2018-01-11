package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	urlparams := mux.Vars(r)
	log.Println(r, urlparams)
	w.Write([]byte(urlparams["hub.challenge"]))
}

func receiveMsg(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	log.Println("post data:", r.Body)
}
