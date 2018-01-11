package main

import (
	"sync"
	"time"
)

type UserState struct {
	state          string
	timeoutChannel chan bool
	stateLock      *sync.Mutex
}

var userStates map[string]UserState

func init() {
	userStates = make(map[string]UserState)
}
func setUserState(userID string, newState string) {
	if _, ok := userStates[userID]; !ok {
		userStates[userID] = UserState{"base", make(chan bool), &sync.Mutex{}}
	}
	currentState := userStates[userID]
	currentState.stateLock.Lock()
	currentState.state = newState
	currentState.stateLock.Unlock()
	userStates[userID] = currentState
	if newState != "base" {
		go func() {
			select {
			case _ = <-currentState.timeoutChannel:
				setUserState(userID, "base")
			case <-time.After(time.Minute * 1):
				sendMsg(userID, "Ok, nevermind (timeout). What would you like to do now?", baseButtons)
				setUserState(userID, "base")
			}
		}()
	}
}
