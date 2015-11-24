package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s cmds...\n", os.Args[0])
		os.Exit(4)
	}

	cmds := exec.Command(os.Args[1])
	cmds.Args = os.Args[2:]
	stime := time.Now()
	cmds.Run()
	etime := time.Now()
	fmt.Fprintf(os.Stdout, "run (%v) time %v\n", os.Args[1:], etime.Sub(stime))
}
