package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var nextId int
var tasksByName = make(map[string]*TaskInfo)
var tasksByID = make(map[int]*TaskInfo)

type Tasker struct {
}

func (t Tasker) AddTask(task *Task, reply *int) error {
	taskInfo := &TaskInfo{ID: nextId, Task: task, Done: false, ExitCode: -1}

	log.Println("AddTask", taskInfo.ID, taskInfo.Task.Name)
	task.Print()

	tasksByName[task.Name] = taskInfo
	tasksByID[taskInfo.ID] = taskInfo

	nextId++

	*reply = taskInfo.ID

	return nil
}

func (t Tasker) EndTask(taskCompletion *TaskCompletion, reply *int) error {
	taskInfo := tasksByID[taskCompletion.ID]

	log.Println("EndTask", taskInfo.Task.Name, taskCompletion.ExitCode)

	taskInfo.ExitCode = taskCompletion.ExitCode
	taskInfo.Done = true

	return nil
}

func (t Tasker) ClearTask(name *string, reply *int) error {
	//taskInfo, present := tasksByName[*name]
	_, present := tasksByName[*name]
	if present {
		delete(tasksByName, *name)
		//delete(tasksByID, taskInfo.ID)
	}

	return nil
}

func (t Tasker) ListTask(name *string, reply *[]TaskInfo) error {
	if *name == "" {
		*reply = make([]TaskInfo, len(tasksByName))
		i := 0
		for _, taskInfo := range tasksByName {
			(*reply)[i] = *taskInfo
			i++
		}
	} else {
		taskInfo, ok := tasksByName[*name]
		if ok {
			*reply = []TaskInfo{*taskInfo}
		} else {
			*reply = []TaskInfo{}
		}
	}
	return nil
}

func startServer() {
	if err := os.RemoveAll(GlobalTaskerConfig.RPC.Socket); err != nil {
		log.Fatal(err)
	}

	t := Tasker{}
	rpc.Register(t)
	rpc.HandleHTTP()
	l, e := net.Listen("unix", GlobalTaskerConfig.RPC.Socket)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Println("Serving...")
	http.Serve(l, nil)
}

