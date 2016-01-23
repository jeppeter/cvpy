package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type IntData struct {
	inner int
}

func (i *IntData) Stringer() string {
	return fmt.Sprintf("%d", i.inner)
}

func (i *IntData) TypeName() string {
	return "IntData"
}

func (i *IntData) Less(j RBTreeData) bool {
	var jv *IntData
	if i.TypeName() != j.TypeName() {
		panic(fmt.Sprintf("i (%s) not type j (%s)", i.TypeName(), j.TypeName()))
	}
	jv = ((*IntData)(unsafe.Pointer((reflect.ValueOf(j).Pointer()))))
	if i.inner < jv.inner {
		return true
	}
	return false
}

func (i *IntData) Equal(j RBTreeData) bool {
	var jv *IntData
	if i.TypeName() != j.TypeName() {
		panic(fmt.Sprintf("i (%s) not type j (%s)", i.TypeName(), j.TypeName()))
	}
	jv = ((*IntData)(unsafe.Pointer((reflect.ValueOf(j).Pointer()))))
	if i.inner == jv.inner {
		return true
	}
	return false
}

func NewIntData(i int) *IntData {
	p := &IntData{}
	p.inner = i
	return p
}

func main() {
	var getdata RBTreeData
	var pi *IntData
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s num\n", os.Args[0])
		os.Exit(4)
	}

	num, _ := strconv.Atoi(os.Args[1])
	nums := []*IntData{}
	getnums := []*IntData{}
	rand.Seed(int64(time.Now().Nanosecond()))
	rbtree := NewRBTree()
	for i := 0; i < num; i++ {
		pi = NewIntData(rand.Int() % (num * 100))
		nums = append(nums, pi)
		rbtree.Insert(pi)
	}
	if true {
		for i := 0; i < num; i++ {
			getdata = rbtree.GetMin()
			if getdata == nil {
				break
			}
			pi = ((*IntData)(unsafe.Pointer((reflect.ValueOf(getdata).Pointer()))))
			getnums = append(getnums, pi)
		}

		for i := 0; i < num; i++ {
			fmt.Fprintf(os.Stdout, "[%d]=(%s) (%s)\n", i, nums[i].Stringer(), getnums[i].Stringer())
		}
	}

	return

}
