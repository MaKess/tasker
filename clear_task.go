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

func parseClearArgs(name *string, sockAddr *string) {

	// TODO: add a real argument parsing here as well

	if len(os.Args) != 3 {
		usageClear()
	}

	*name = os.Args[2]
	*sockAddr = SockAddr
}

func clearTask() {
	var name string
	var sockAddr string

	parseClearArgs(&name, &sockAddr)

	client, err := rpc.DialHTTP("unix", sockAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply int
	err = client.Call("Tasker.ClearTask", name, &reply)
	if err != nil {
		log.Fatal("error while reporting task:", err)
	}
}
