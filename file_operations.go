package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

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
	bytes, err := json.MarshalIndent(record, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile("record.json", bytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func saveRecords(records []FileRecord) error {
	jsonBytes, err := json.MarshalIndent(records, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile("batch_record.json", jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func processSingleFile(filepathArg string) error {
	info, err := os.Stat(filepathArg)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found")
		}
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("that's a directory not a file. Use --batch to process directories")
	}

	fmt.Println("Name:", info.Name())
	fmt.Println("Size:", info.Size())
	fmt.Println("Modified:", info.ModTime().Format("2006-01-02 15:04:05"))

	hash, err := hashFile(filepathArg)
	if err != nil {
		return fmt.Errorf("error hashing file: %v", err)
	}
	fmt.Println("SHA-256:", hash)

	record := FileRecord{
		Filename:   info.Name(),
		Size:       info.Size(),
		Modified:   info.ModTime().Format("2006-01-02 15:04:05"),
		SHA256:     hash,
		RecordedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = saveRecord(record)
	if err != nil {
		return fmt.Errorf("error saving record: %v", err)
	}
	fmt.Println("Record saved to record.json")
	return nil
}

func batchHash(dirPath string) error {
	var records []FileRecord

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		hash, err := hashFile(path)
		if err != nil {
			return err
		}
		record := FileRecord{
			Filename:   info.Name(),
			Size:       info.Size(),
			Modified:   info.ModTime().Format("2006-01-02 15:04:05"),
			SHA256:     hash,
			RecordedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		records = append(records, record)
		return nil
	})
	if err != nil {
		return err
	}

	err = saveRecords(records)
	if err != nil {
		return err
	}

	fmt.Printf("Processed %d files. Record saved to batch_record.json\n", len(records))
	return nil
}
