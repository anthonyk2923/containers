package main

import (
	"log"
	"syscall"
)

func Namespace() {
	err := syscall.Unshare(
		syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET,
	)
	if err != nil {
		log.Fatal("error setting up namespace: ", err)
	}
}

