package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"time"
)

func parseAddArgs(task *Task, sockAddr *string) {
	base := fmt.Sprintf("%s add", os.Args[0])
	addFlags := flag.NewFlagSet(base, flag.ExitOnError)
	addFlags.StringVar(&task.Name, "n", "", "the name of this task")
	dep := addFlags.String("d", "", "one or more comma separated tasks on which this task depends")
	addFlags.StringVar(&task.Checker, "c", "", "an optional check to be run whose return value is used instead")
	//addFlags.StringVar(sockAddr, "s", SockAddr, "path to the socket used for coordination")
	*sockAddr = SockAddr
	addFlags.Parse(os.Args[2:])

	if *dep != "" {
		task.Dep = strings.Split(*dep, ",")
	}

	task.Cmd = addFlags.Args()

	ok := true

	if task.Name == "" {
		fmt.Println("need to provide task name '-n'")
		ok = false
	}

	if !ok {
		fmt.Println()
		addFlags.Usage()
		os.Exit(1)
	}

	task.Dir = getcwd()
	task.Env = os.Environ()
}

func addTaskEnd(client *rpc.Client, id int, exitCode int) {
	var reply int
	err := client.Call("Tasker.EndTask", &TaskCompletion{ID: id, ExitCode: exitCode}, &reply)
	if err != nil {
		log.Fatal("error while reporting task:", err)
	}
	log.Println("task finished with code", exitCode)
	os.Exit(exitCode)
}

func addTask() {
	var sockAddr string
	task := Task{}

	parseAddArgs(&task, &sockAddr)

	task.Print()

	client, err := rpc.DialHTTP("unix", sockAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var id int
	err = client.Call("Tasker.AddTask", &task, &id)
	if err != nil {
		log.Fatal("error while adding task:", err)
	}
	//log.Printf("new task '%s' has id '%d'\n", task.Name, id)


	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for range c {
			fmt.Println() // empy line in case there's someting else printed
			log.Println("caught Ctrl+C")
			addTaskEnd(client, id, 130) // exit code 130 for ^C
		}
	}()


	if len(task.Dep) > 0 {
		log.Println("waiting for dependencies:", task.Dep)

		start := time.Now()

		for missingDep := true; missingDep; {

			fmt.Printf("\033[2K\033[1Gwaiting for %v", time.Since(start))

			missingDep = false
			for _, element := range task.Dep {
				var taskInfos []TaskInfo
				err = client.Call("Tasker.ListTask", &element, &taskInfos)
				if err != nil {
					log.Fatal(err)
				} else if len(taskInfos) < 1 {
					missingDep = true
				} else if !taskInfos[0].Done {
					missingDep = true
				} else if taskInfos[0].ExitCode > 0 {
					log.Fatal("a dependency has failed:", element)
				}
			}

			if missingDep {
				time.Sleep(time.Second)
			}
		}

		fmt.Println()
	}

	exitCode := launchTask(&task)

	if task.Checker != "" {
		exitCode = runCmd(makeCmd([]string{task.Checker}))
	}

	addTaskEnd(client, id, exitCode)
}
