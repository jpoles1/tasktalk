package main

import "strconv"

var addTaskButton = ReplyButton{"text", "Add Task", "addTask", ""}
var getTasksButton = ReplyButton{"text", "Get Tasks", "getTasks", ""}
var deleteTaskButton = ReplyButton{"text", "Delete Task", "addTask", ""}
var cancelButton = ReplyButton{"text", "Cancel", "cancel", ""}

var baseButtons = []ReplyButton{addTaskButton, getTasksButton, deleteTaskButton}

func cancelResponse(recipientID string) {
	setUserState(recipientID, "base")
	userStates[recipientID].timeoutChannel <- true
	sendMsg(recipientID, "Ok, nevermind. What would you like to do now?", baseButtons)
}

func formatTaskList(taskList []UserTask) string {
	taskText := "Task List:\n"
	for index, task := range taskList {
		taskText += strconv.Itoa(index+1) + ") - " + task.TaskText + "\n"
	}
	return taskText
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
	sendMsg(recipientID, "Ok, adding your task:\n"+msgText, baseButtons)
	getTasksResponse(recipientID, baseButtons)
}

//Deleting tasks
func deletingTaskResponse(recipientID string) {
	setUserState(recipientID, "deleteTask")
	getTasksResponse(recipientID, []ReplyButton{})
	sendMsg(recipientID, "What task should I delete (enter the task #)?", []ReplyButton{cancelButton})
}
func deletedTaskResponse(recipientID string, msgText string) {
	if msgIndex, err := strconv.Atoi(msgText); err == nil {
		err := dbDeleteTask(recipientID, msgIndex)
		if err == nil {
			setUserState(recipientID, "base")
			sendMsg(recipientID, "Ok, deleting task #"+msgText, baseButtons)
			getTasksResponse(recipientID, baseButtons)
		} else {
			sendMsg(recipientID, "Error: invalid task #!", []ReplyButton{cancelButton})
		}
	} else {
		sendMsg(recipientID, "Error: response was not a #!", []ReplyButton{cancelButton})
	}
}
