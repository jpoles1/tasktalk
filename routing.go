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
	Object string `json:"object"`
	Entry  []struct {
		Messaging []struct {
			Message struct {
				Text string `json:"text"`
				Seq  int    `json:"seq"`
				Mid  string `json:"mid"`
			} `json:"message"`
			Timestamp int64 `json:"timestamp"`
			Sender    struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
		} `json:"messaging"`
		Time int64  `json:"time"`
		ID   string `json:"id"`
	} `json:"entry"`
}

func receiveMsg(w http.ResponseWriter, r *http.Request) {
	var postData IncomingMessage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postData)
	if err != nil {
		panic(err)
	}
	if len(postData.Entry) < 1 {
		log.Println("Error: Malformed Request")
		return
	}
	if len(postData.Entry[0].Messaging) < 1 {
		log.Println("Error: Malformed Request")
		return
	}
	log.Println("Message Data:", postData)
	msgText := postData.Entry[0].Messaging[0].Message.Text
	w.Write([]byte("Echo: " + msgText))
}
