package main

import (
	"bufio"
	"fmt"
	"github.com/u-root/u-root/pkg/ldd"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Setup(reqPath, sandboxPath string) {
	file, err := os.Open(reqPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	os.MkdirAll(filepath.Join(sandboxPath, "lib"), 0755)
	os.MkdirAll(filepath.Join(sandboxPath, "lib64"), 0755)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		line := scanner.Text()
		destPath := filepath.Join(sandboxPath, filepath.Base(line))
		if err := copy(line, destPath); err != nil {
			log.Panic(err)
		}
		if err := copyLibs(line, sandboxPath); err != nil {
			log.Panic("failed to copy libs: ", err)
		}
		if err := os.Chmod(destPath, 0755); err != nil {
			log.Panicf("Failed to set execute permissions on %s: %v", destPath, err)
		}
	}

	if err := applyChmodRecursive(sandboxPath, 0755); err != nil {
		log.Panicf("Failed to apply chmod recursively: %v", err)
	}
}

func applyChmodRecursive(rootPath string, mode os.FileMode) error {
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err := os.Chmod(path, mode); err != nil {
			return fmt.Errorf("failed to chmod %s: %v", path, err)
		}
		return nil
	})
	return err
}
func copyLibs(binary, sandboxPath string) error {
	libraries, err := ldd.List(binary)
	if err != nil {
		log.Fatalf("Error retrieving libraries for %s: %v", binary, err)
	}

	for _, lib := range libraries {
		destPath := filepath.Join(sandboxPath, lib)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			log.Printf("Failed to create directory for %s: %v", lib, err)
			continue
		}

		if err := copy(lib, destPath); err != nil {
			log.Printf("Failed to copy %s: %v", lib, err)
			continue
		}
	}

	return nil
}

func copy(fromPath, toPath string) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
