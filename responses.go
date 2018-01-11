package main

import "strconv"

var addTaskButton = ReplyButton{"text", "Add Task", "addTask", ""}
var getTasksButton = ReplyButton{"text", "Get Tasks", "getTasks", ""}
var deleteTaskButton = ReplyButton{"text", "Delete Task", "addTask", ""}
var cancelButton = ReplyButton{"text", "Cancel", "cancel", ""}

var baseButtons = []ReplyButton{addTaskButton, getTasksButton, deleteTaskButton}

func cancelResponse(recipientID string) {
	setUserState(recipientID, "base")
	sendMsg(recipientID, "Ok, nevermind. What would you like to do now?", baseButtons)
}

func getTasksResponse(recipientID string, buttonList []ReplyButton) {
	sendMsg(recipientID, formatTaskList(dbFetchTasks(recipientID)), buttonList)
}

//Adding tasks
func addingTaskResponse(recipientID string) {
	setUserState(recipientID, "addTask")
	sendMsg(recipientID, "What task can I add to your list?", []ReplyButton{cancelButton})
}
func addedTaskResponse(recipientID string, msgText string) {
	setUserState(recipientID, "base")
	dbAddTask(recipientID, msgText)
	sendMsg(recipientID, "Ok, adding your task: "+msgText, baseButtons)
	getTasksResponse(recipientID, []ReplyButton{})
}

//Deleting tasks
func deletingTaskResponse(recipientID string) {
	setUserState(recipientID, "deleteTask")
	sendMsg(recipientID, "What task should I delete (enter the task #)?", []ReplyButton{cancelButton})
	getTasksResponse(recipientID, []ReplyButton{})
}
func deletedTaskResponse(recipientID string, msgText string) {
	if msgIndex, err := strconv.Atoi(msgText); err == nil {
		err := dbDeleteTask(recipientID, msgIndex)
		if err == nil {
			setUserState(recipientID, "base")
			sendMsg(recipientID, "Ok, deleting task #"+msgText, baseButtons)
		} else {
			sendMsg(recipientID, "Error: invalid task #!", baseButtons)
		}
	} else {
		sendMsg(recipientID, "Error: response was not a #!", []ReplyButton{cancelButton})
	}
}
