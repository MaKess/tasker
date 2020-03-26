package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type TaskerRPCConfig struct {
	Socket string `json:"socket"`
}

type TaskerConfig struct {
	RPC TaskerRPCConfig `json:"rpc"`
}

// default values for the configuration
var GlobalTaskerConfig = TaskerConfig {
	RPC: TaskerRPCConfig {
		Socket: "/tmp/tasker-%s.sock",
	},
}

func readConfig() {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	for _, fileName := range [...]string{"/etc/tasker.json", "config.json", "~/.tasker.json"} {
		if strings.HasPrefix(fileName, "~/") {
			fileName = filepath.Join(homeDir, fileName[2:])
		}

		configFile, err := os.Open(fileName)
		if err == nil {
			json.NewDecoder(configFile).Decode(&GlobalTaskerConfig)
			configFile.Close()
		}
	}

	// fix up dynamic configuration items
	GlobalTaskerConfig.RPC.Socket = fmt.Sprintf(GlobalTaskerConfig.RPC.Socket, usr.Username)
}

type Task struct {
	Name string
	Dep []string
	Checker string
	Env []string
	Dir string
	Cmd []string
}

func (task *Task) Print() {
	log.Println("name:", task.Name)
	log.Println("dep:", task.Dep)
	log.Println("checker:", task.Checker)
	//log.Println("env:", task.Env)
	log.Println("dir:", task.Dir)
	log.Println("cmd:", task.Cmd)
}

type TaskCompletion struct {
	ID int
	ExitCode int
}

type TaskInfo struct {
	ID int
	Done bool
	Task *Task
	ExitCode int
}

func getcwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("could not reliably determine current directory", err)
		cwd = ""
	}
	return cwd
}

func makeCmd(args []string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

func runCmd(cmd *exec.Cmd) int {
	if err := cmd.Run() ; err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			return -1
		}
	} else {
		return 0
	}
}

func launchTask(task *Task) int {
	log.Println("launching task...")
	cmd := makeCmd(task.Cmd)
	cmd.Dir = task.Dir
	cmd.Env = task.Env
	return runCmd(cmd)
}

func usage() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s <command>\n", os.Args[0])
	fmt.Println("command:")
	fmt.Println("\tadd")
	fmt.Println("\tclear")
	fmt.Println("\tcheck")
	fmt.Println("\tlist")
	fmt.Println("\tserver")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	readConfig()

	switch os.Args[1] {
	case "add":
		addTask()
	case "clear":
		clearTask()
	case "check":
		checkTask()
	case "list":
		listTask()
	case "server":
		startServer()
	default:
		usage()
	}
}
