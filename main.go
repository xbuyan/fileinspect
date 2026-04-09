package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: fileinspect <filepath>")
		return
	}

	filepath := os.Args[1]

	info, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Error: file not found")
		} else {
			fmt.Println("Error", err)
		}

		return
	}
	if info.IsDir() {
		fmt.Println("Error: that's a directory not a file")
		return
	}

	fmt.Println("Name", info.Name())
	fmt.Println("Size", info.Size())
	fmt.Println("Modified:", info.ModTime().Format("2006-01-02 15:04:05"))

	hash, err := hashFile(filepath)
	if err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}
	fmt.Println("SHA-256: ", hash)
}

func hashFile(filepath string) (string, error) {
	content, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer content.Close()
	hasher := sha256.New()

	io.Copy(hasher, content)
	bytes := hasher.Sum(nil)

	return hex.EncodeToString(bytes), nil
}
