package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	// Print current process ID to show we can access the host's PID namespace
	fmt.Println("Current Process ID (PID):", syscall.Getpid())

	// Attempting to read the /proc directory (which lists processes on the host system)
	// This will demonstrate that we can still see processes running on the host despite being inside the chroot jail.
	procDir := "/proc"
	entries, err := os.ReadDir(procDir)
	if err != nil {
		log.Fatalf("Failed to read /proc: %v", err)
	}

	// List process entries (This should list processes on the host, showing lack of PID isolation)
	fmt.Println("Processes visible from /proc (should be global processes from host system):")
	for _, entry := range entries {
		fmt.Println(entry.Name()) // Displays process information for processes outside the chroot
	}

	// Try to access a file outside the chroot to show it’s not properly isolated
	// This file should exist on the host system, showing that chroot doesn't fully isolate the filesystem.
	file, err := os.Open("/etc/hostname") // Trying to access a file outside the chroot
	if err != nil {
		log.Printf("Expected error trying to access file outside the chroot: %v", err)
	} else {
		fmt.Println("Successfully accessed /etc/hostname outside of chroot:", file.Name())
		defer file.Close()
	}
}
