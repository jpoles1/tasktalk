package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
type ReplyButton struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
	ImageURL    string `json:"image_url"`
}

var cancelButton ReplyButton = ReplyButton{"text", "Cancel", "cancel", ""}
var addTaskButton ReplyButton = ReplyButton{"text", "Add Task", "addTask", ""}

type OutgoingMessage struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text         string        `json:"text"`
		ReplyButtons []ReplyButton `json:"quick_replies"`
	} `json:"message"`
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
	senderID := postData.Entry[0].Messaging[0].Sender.ID
	msgText := postData.Entry[0].Messaging[0].Message.Text
	if val, ok := userStates[senderID]; ok && val.state != "base" && msgText != "" {
		if val.state == "addTask" {
			if msgText == "Cancel" {
				setUserState(senderID, "base")
				sendMsg(senderID, "Ok, nevermind. What would you like to do now?", []ReplyButton{addTaskButton})
			} else {
				setUserState(senderID, "base")
				sendMsg(senderID, "Ok, adding your task: "+msgText, []ReplyButton{addTaskButton})
			}
		}
	} else if msgText == "Add Task" {
		setUserState(senderID, "addTask")
		sendMsg(senderID, "What task can I add to your list?", []ReplyButton{cancelButton})
	} else if msgText == "Cancel" {
		setUserState(senderID, "base")
		sendMsg(senderID, "Ok, nevermind. What would you like to do now?", []ReplyButton{addTaskButton})
	} else if msgText != "" {
		msgText = "Echo: " + msgText
		//Quick reply buttons
		sendMsg(senderID, msgText, []ReplyButton{addTaskButton})
	} else {
		var postData interface{}
		log.Println("Unknown Message Format!")
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&postData)
		if err != nil {
			panic(err)
		}
		log.Println("Raw Data:", r.Body)
		log.Println("Message Data:", postData)
	}
	w.Write([]byte("ok"))
}

func sendJSON(jsonData []byte) {
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + fbToken
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))
}
func sendMsg(recipientID string, msgText string, replyButtons []ReplyButton) {
	msgData := OutgoingMessage{}
	msgData.Recipient.ID = recipientID
	msgData.Message.Text = msgText
	msgData.Message.ReplyButtons = replyButtons
	jsonData, _ := json.Marshal(msgData)
	sendJSON(jsonData)
}
