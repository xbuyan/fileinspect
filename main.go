package main

import (
	"flag"
	"fmt"
	"os"
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
		if err := verifyFile(filepathArg); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return
	}

	// Handle --batch flag second
	if *batch {
		if err := batchHash(filepathArg); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return
	}

	// If no flags, treat as single file mode
	if err := processSingleFile(filepathArg); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
