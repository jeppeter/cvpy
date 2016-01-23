package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

var __syscall_outputdebugstring *syscall.Proc

func init() {
	__syscall_outputdebugstring = nil
	d, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return
	}
	__syscall_outputdebugstring, err = d.FindProc("OutputDebugStringW")
	if err != nil {
		__syscall_outputdebugstring = nil
		return
	}
	return
}

func InnerDebugOutput(s string) error {
	if __syscall_outputdebugstring != nil {
		p := syscall.StringToUTF16Ptr(s)
		__syscall_outputdebugstring.Call(uintptr(unsafe.Pointer(p)))
	}
	return nil
}

func Debug(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stdout, s)
	InnerDebugOutput(s)
	return len(s)
}

func Error(format string, a ...interface{}) int {
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("[%s:%d]\t", f, l)
	s += fmt.Sprintf(format, a...)
	s += "\n"
	fmt.Fprint(os.Stderr, s)
	InnerDebugOutput(s)
	return len(s)
}
