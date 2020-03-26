package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func usageClear() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s clear <name>\n", os.Args[0])
	os.Exit(1)
}

func clearTask() {
	if len(os.Args) != 3 {
		usageClear()
		return
	}

	name := os.Args[2]

	client, err := rpc.DialHTTP("unix", GlobalTaskerConfig.RPC.Socket)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply int
	err = client.Call("Tasker.ClearTask", name, &reply)
	if err != nil {
		log.Fatal("error while reporting task:", err)
	}
}
