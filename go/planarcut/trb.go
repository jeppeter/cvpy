package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type IntData int

func (i *IntData) Stringer() string {
	return fmt.Sprintf("%d", *i)
}

func (i *IntData) Less(j *IntData) bool {
	if *i < *j {
		return true
	}
	return false
}

func (i *IntData) Equal(j *IntData) bool {
	if *i == *j {
		return true
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s num\n", os.Args[0])
		os.Exit(4)
	}

	num, _ := strconv.Atoi(os.Args[1])
	nums := []IntData{}
	getnums := []IntData{}
	rand.Seed(float64(time.Now().Nanosecond()))
	rbtree := NewRBTree()
	for i := 0; i < num; i++ {
		ni := rand.Int31n((num * 100))
		nums = append(nums, ni)
		rbtree.Insert(ni)
	}

	for i := 0; i < num; i++ {
		ni := rbtree.GetMin()
		getnums = append(getnums, ni.Data)
	}

	fmt.Fprintf(os.Stdout, "random get (%v)\n", nums)
	fmt.Fprintf(os.Stdout, "sort by rbtree (%v)\n", getnums)
	return

}
