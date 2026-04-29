package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func verifyFile(filepath string) error {
	data, err := os.ReadFile("record.json")
	if err != nil {
		return fmt.Errorf("cannot read record.json: %v", err)
	}

	var record FileRecord
	err = json.Unmarshal(data, &record)
	if err != nil {
		return fmt.Errorf("cannot parse record.json: %v", err)
	}

	currentHash, err := hashFile(filepath)
	if err != nil {
		return fmt.Errorf("cannot hash file: %v", err)
	}

	if currentHash == record.SHA256 {
		fmt.Println("Recorded hash:", record.SHA256)
		fmt.Println("Current hash:", currentHash)
		fmt.Println("VERIFIED: File has not been tampered with.")
	} else {
		fmt.Println("Recorded hash:", record.SHA256)
		fmt.Println("Current hash:", currentHash)
		fmt.Println("ALERT: File has been modified since recording.")
	}
	return nil
}
