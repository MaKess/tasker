package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func usageList() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s list\n", os.Args[0])
	os.Exit(1)
}

func listTask() {
	if len(os.Args) != 2 {
		usageList()
		return
	}

	client, err := rpc.DialHTTP("unix", GlobalTaskerConfig.RPC.Socket)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var taskInfos []TaskInfo
	noName := ""
	err = client.Call("Tasker.ListTask", &noName, &taskInfos)
	if err != nil {
		log.Fatal(err)
	}

	if len(taskInfos) > 0 {
		log.Println("tasks:")
		for _, taskInfo := range taskInfos {
			status := "?"
			if !taskInfo.Done {
				status = "not done"
			} else if taskInfo.ExitCode == 0 {
				status = "done"
			} else if taskInfo.ExitCode > 0 {
				status = "failed"
			}
			log.Printf(" - %s: %s, depends: %v, command: %v",
				taskInfo.Task.Name,
				status,
				taskInfo.Task.Dep,
				taskInfo.Task.Cmd)
		}
	} else {
		log.Println("there are no queded tasks")
	}
}
