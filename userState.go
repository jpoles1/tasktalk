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
	go func() {
		select {
		case _ = <-currentState.timeoutChannel:
			sendMsg(userID, "Ok, nevermind. What would you like to do now?", []ReplyButton{addTaskButton})
		case <-time.After(time.Minute * 1):
			setUserState(userID, "base")
		}
	}()
}
