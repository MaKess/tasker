package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func usageCheck() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s check <name>\n", os.Args[0])
	os.Exit(1)
}

func checkTask() {
	if len(os.Args) != 3 {
		usageCheck()
		return
	}

	name := os.Args[2]

	client, err := rpc.DialHTTP("unix", GlobalTaskerConfig.RPC.Socket)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var taskInfos []TaskInfo
	err = client.Call("Tasker.ListTask", &name, &taskInfos)
	if err != nil {
		log.Fatal(err)
	} else if len(taskInfos) > 0 {
		taskInfo := &taskInfos[0]
		taskInfo.Task.Print()
		if taskInfo.Done {
			log.Println("task finished with code", taskInfo.ExitCode)
		} else {
			log.Println("task has not finished yet")
		}
	} else {
		log.Fatalf("task '%s' not found\n", name)
	}
}
