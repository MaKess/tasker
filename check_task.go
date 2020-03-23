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

func parseCheckArgs(name *string, sockAddr *string) {

	// TODO: add a real argument parsing here as well

	if len(os.Args) != 3 {
		usageCheck()
	}

	*name = os.Args[2]
	*sockAddr = SockAddr
}

func checkTask() {
	var name string
	var sockAddr string

	parseCheckArgs(&name, &sockAddr)

	client, err := rpc.DialHTTP("unix", sockAddr)
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
