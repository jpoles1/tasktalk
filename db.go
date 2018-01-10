package main

import (
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
)

var mongoURI string
var mongoDB *mgo.Session

func mongoLoad() {
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
