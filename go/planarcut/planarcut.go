package main

import (
	"fmt"
	"os"
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
	planar, err := MakePlanarGraph(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse %s error(%s)\n", os.Args[1], err.Error())
		os.Exit(5)
	}
	planar.DebugGraph()
}
