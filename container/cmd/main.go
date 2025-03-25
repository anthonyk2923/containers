package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run ./cmd <command>\n")
	}
	reqPath := "../../req.binaries"
	sandboxPath := "../sandbox"

	binaryPath := os.Args[1]
	args := os.Args[2:]

	Namespace()
	Setup(reqPath, sandboxPath)
	Chroot(sandboxPath)

	fmt.Println("Attempting to execute:", binaryPath, args)
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing: %v", err)
	}
}
