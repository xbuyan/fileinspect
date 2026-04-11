package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type FileRecord struct {
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	Modified   string `json:"modified"`
	SHA256     string `json:"sha256"`
	RecordedAt string `json:"recorded_at"`
}

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

	record := FileRecord{
		Filename:   info.Name(),
		Size:       info.Size(),
		Modified:   info.ModTime().Format("2006-01-02 15:04:05"),
		SHA256:     hash,
		RecordedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = saveRecord(record)
	if err != nil {
		fmt.Println("Error saving record:", err)
		return
	}
	fmt.Println("Record saved to record.json")
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

func saveRecord(record FileRecord) error {

	// 1. Convert the struct to JSON with json.MarshalIndent

	bytes, err := json.MarshalIndent(record, "", "\t")

	if err != nil {

		return err
	}
	err = os.WriteFile("record.json", bytes, 0644)

	if err != nil {

		return err
	}
	return nil

}
