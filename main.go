package main

import (
	"fmt"
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

		fmt.Println("Error", err)
		return
	}

	fmt.Println("Name", info.Name())
	fmt.Println("Size", info.Size())
	fmt.Println("Modified:", info.ModTime().Format("2006-01-02 15:04:05"))
}
