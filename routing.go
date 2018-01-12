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

//ReplyButton contains the data for a quick-reply button to be included in an OutgoingMessage
type ReplyButton struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
	ImageURL    string `json:"image_url"`
}

//OutgoingMessage contains the data included in a facebook message
type OutgoingMessage struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text         string        `json:"text"`
		ReplyButtons []ReplyButton `json:"quick_replies,omitempty"`
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
		if msgText == "Cancel" {
			cancelResponse(senderID)
		} else if val.state == "addTask" {
			addedTaskResponse(senderID, msgText)
		} else if val.state == "deleteTask" {
			deletedTaskResponse(senderID, msgText)
		}
	} else if msgText == "Add Task" {
		addingTaskResponse(senderID)
	} else if msgText == "Delete Task" {
		deletingTaskResponse(senderID)
	} else if msgText == "Get Tasks" {
		getTasksResponse(senderID, baseButtons)
	} else if msgText == "Cancel" {
		cancelResponse(senderID)
	} else if msgText != "" {
		msgText = "Echo: " + msgText
		//Quick reply buttons
		sendMsg(senderID, msgText, baseButtons)
	} else {
		log.Println("Unknown Message Format!")
		log.Println("Raw Data:", r.Body)
	}
	w.Write([]byte("ok"))
}

func sendJSON(url string, jsonData []byte) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + fbToken
	sendJSON(url, jsonData)
}
