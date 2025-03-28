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

var copiedLibs = make(map[string]bool)

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
		if err := copy(line, destPath, sandboxPath); err != nil {
			log.Panic(err)
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
		log.Printf("Error retrieving libraries for %s: %v", binary, err)
		return nil
	}

	for _, lib := range libraries {
		if copiedLibs[lib] {
			continue
		}

		copiedLibs[lib] = true
		destPath := filepath.Join(sandboxPath, lib)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			log.Printf("Failed to create directory for %s: %v", lib, err)
			continue
		}

		if err := copy(lib, destPath, sandboxPath); err != nil {
			log.Printf("Failed to copy %s: %v", lib, err)
			continue
		}
	}

	return nil
}

func copy(fromPath, toPath, sandboxPath string) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		files, err := os.ReadDir(fromPath)
		if err != nil {
			return fmt.Errorf("failed to read directory: %v", err)
		}

		for _, file := range files {
			if !file.IsDir() {
				filePath := filepath.Join(fromPath, file.Name())
				destPath := filepath.Join(sandboxPath, filepath.Base(filePath))
				if err := copy(filePath, destPath, sandboxPath); err != nil {
					log.Printf("Failed to copy %s: %v", filePath, err)
					continue
				}

				if err := copyLibs(filePath, sandboxPath); err != nil {
					log.Printf("Failed to copy libs for %s: %v", filePath, err)
					continue
				}

				if err := os.Chmod(destPath, 0755); err != nil {
					log.Printf("Failed to set execute permissions on %s: %v", destPath, err)
					continue
				}
			}
		}
		return nil
	}
	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	if err := copyLibs(fromPath, sandboxPath); err != nil {
		return fmt.Errorf("failed to copy libs: %v", err)
	}

	if err := os.Chmod(toPath, 0755); err != nil {
		return fmt.Errorf("Failed to set execute permissions on %s: %v", toPath, err)
	}

	return nil
}
