package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	urlparams := mux.Vars(r)
	log.Println(r, urlparams)
	w.Write([]byte("Hello World!"))
}
