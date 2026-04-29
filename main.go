package main

import (
	"flag"
	"fmt"
)

func main() {
	batch := flag.Bool("batch", false, "hash all files in a directory")
	verify := flag.Bool("verify", false, "verify file against saved record")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: fileinspect [--batch] [--verify] <path>")
		return
	}

	filepathArg := flag.Arg(0)

	// Handle --verify flag first
	if *verify {
		err := verifyFile(filepathArg)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	// Handle --batch flag second
	if *batch {
		err := batchHash(filepathArg)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	// If no flags, treat as single file mode
	err := processSingleFile(filepathArg)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
