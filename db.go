package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserTask struct {
	TaskID   bson.ObjectId `bson:"_id"`
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
	newTask := UserTask{bson.NewObjectId(), userID, taskText}
	mongoSesh := mongoDB.Copy()
	defer mongoSesh.Close()
	err := mongoSesh.DB("heroku_r47fhcrt").C("tasks").Insert(newTask)
	if err != nil {
		fmt.Println("Failure to insert task document:\n", err)
	}
	userStates[userID].timeoutChannel <- true
}
func dbDeleteTask(userID string, taskIndex int) error {
	mongoSesh := mongoDB.Copy()
	defer mongoSesh.Close()
	taskList := dbFetchTasks(userID)
	if len(taskList) < taskIndex {
		currentTask := taskList[taskIndex]
		err := mongoSesh.DB("heroku_r47fhcrt").C("tasks").Remove(bson.M{"_id": currentTask.TaskID})
		if err != nil {
			fmt.Println("Failure to insert account document:\n", err)
		}
		userStates[userID].timeoutChannel <- true
		return nil
	}
	return errors.New("Invalid Task Index!")
}
func dbFetchTasks(userID string) []UserTask {
	mongoSesh := mongoDB.Copy()
	defer mongoSesh.Close()
	var taskList []UserTask
	mongoSesh.DB("heroku_r47fhcrt").C("tasks").Find(bson.M{"userid": userID}).All(&taskList)
	return taskList
}
