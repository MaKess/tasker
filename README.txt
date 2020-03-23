compile with
$ go build tasker.go server.go add_task.go clear_task.go check_task.go list_task.go


launch a server in a background terminal:
$ ./tasker server


queue tasks with their corresponding names and dependencies (each to be executed in a separate terminal).
the actual command to be executed is the last part after the "--".

environment 1:
$ ./tasker add -n lulz -d foo,bar -- ls

environment 2:
$ ./tasker add -n bar -- sleep 5

environment 3:
$ ./tasker add -n foo -- date

this will queue three tasks.
the first, we name it "lulz", depends on the successful completion (exit code 0) of tasks "foo" and "bar".
so "tasker" will only launch "ls" once the other two are done.

to monitor what's going on, the following commands are available:
so we can now 

$ ./tasker list
lists all the queued up commands and whether they are done

$ ./tasker clear foo
removes the result of a given task, so we can run it again (otherwise "lulz" will immediately think "foo" would have completed)

$ ./tasker check foo
show more details of the queued task "foo"
