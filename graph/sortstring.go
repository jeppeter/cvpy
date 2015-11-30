package main

import (
	"strings"
)

type StringSort []string

func (a StringSort) Len() int {
	return len(a)
}

func (a StringSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
	return
}

func (a StringSort) Less(i, j int) bool {
	if strings.Compare(a[i], a[j]) < 0 {
		return true
	}
	return false
}
