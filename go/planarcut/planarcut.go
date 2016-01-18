package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Usage(ec int, format string, a ...interface{}) {
	fp := os.Stderr
	if ec == 0 {
		fp = os.Stdout
	}

	if format != "" {
		s := fmt.Sprintf(format, a...)
		s += "\n"
		fmt.Fprint(fp, s)
	}

	fmt.Fprintf(fp, "planarcut infile\n")
	os.Exit(ec)
}

func main() {
	if len(os.Args) < 2 {
		Usage(3, "need one arg")
	}
	log.SetFlags(log.Lshortfile)
	planar, err := MakePlanarGraph(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s error(%s)\n", os.Args[1], err.Error())
		os.Exit(5)
	}
	stime := time.Now()
	maxflow, err := planar.GetMaxFlow()
	etime := time.Now()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not get maxflow err (%s)\n", err.Error())
		os.Exit(6)
	}
	fmt.Fprintf(os.Stdout, "maxflow %f (%s)\n", maxflow, etime.Sub(stime))
	return
}
