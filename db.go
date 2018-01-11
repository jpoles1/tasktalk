package main

import (
	"fmt"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserTask struct {
	UserID   string
	TaskText string
}

var mongoURI string
var mongoDB *mgo.Session

func dbLoad() {
	mongoURI = os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("No MongoDB URI supplied in .env config file!")
	}
	var err error
	mongoDB, err = mgo.Dial(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to provided MongoDB URI:\n", err)
	}
}

func dbAddTask(userID string, taskText string) {
	newTask := UserTask{userID, taskText}
	mongoSesh := mongoDB.Copy()
	defer mongoSesh.Close()
	err := mongoSesh.DB("heroku_r47fhcrt").C("tasks").Insert(newTask)
	if err != nil {
		fmt.Println("Failure to insert account document:\n", err)
	}
}
func dbFetchTasks(userID string) string {
	mongoSesh := mongoDB.Copy()
	defer mongoSesh.Close()
	var taskList []UserTask
	mongoSesh.DB("heroku_r47fhcrt").C("tasks").Find(bson.M{"userid": userID}).All(&taskList)
	return fmt.Sprintln(taskList)
}
