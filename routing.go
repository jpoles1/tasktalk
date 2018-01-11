package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	urlparams := mux.Vars(r)
	log.Println(r, urlparams)
	w.Write([]byte(urlparams["hub.challenge"]))
}

//IncomingMessage contains the data from a facebook message
type IncomingMessage struct {
	entry []struct {
		messaging []struct {
			message struct {
				text string
			}
			timestamp string
			sender    struct {
				id string
			}
			recipient struct {
				id string
			}
		}
	}
}

func receiveMsg(w http.ResponseWriter, r *http.Request) {
	var postData interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postData)
	if err != nil {
		panic(err)
	}

	log.Println("Message Data:", postData)
}
