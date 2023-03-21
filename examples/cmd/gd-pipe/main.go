package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cvvz/go-daemon"
)

func main() {
	r, w, _ := os.Pipe()

	cntxt := &daemon.Context{
		PidFileName: "sample.pid",
		PidFilePerm: 0644,
	}
	cntxt.SetLogFile(w)
	defer cntxt.Release()

	var sigchild chan os.Signal
	if !daemon.WasReborn() {
		sigchild = make(chan os.Signal, 1)
		signal.Notify(sigchild, syscall.SIGCHLD)
	}

	child, err := cntxt.Reborn()
	if err != nil {
		panic(err)
	}

	if child == nil {
		fmt.Printf("from pipe w!\n")
		return
	}

	<-sigchild
	buf := make([]byte, 1024)
	r.Read(buf)
	fmt.Printf("parent receive: %s", buf)

}
