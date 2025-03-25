package main

import (
	"log"
	"os"
	"syscall"
)

func Chroot(sandboxDir string) {
	err := syscall.Chroot(sandboxDir)
	if err != nil {
		log.Fatalf("Chroot failed: %v", err)
	}

	err = os.Chdir("/")
	if err != nil {
		log.Fatalf("Chdir failed: %v", err)
	}

}
