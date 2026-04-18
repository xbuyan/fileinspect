package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
<<<<<<< HEAD
	batch := flag.Bool("batch", false, "hash all files in a directory")
=======

	batch := flag.Bool("batch", false, "hash all files in a directory")

>>>>>>> 6662348 (feat: add --batch flag for directory hashing)
	verify := flag.Bool("verify", false, "verify file against saved record")

	flag.Parse()

	if flag.NArg() < 1 {
<<<<<<< HEAD
		fmt.Println("Usage: fileinspect [--batch] [--verify] <path>")
		return
	}

	filepathArg := flag.Arg(0)
=======

		fmt.Println("Usage: file inspect[--verify] <filepath")
		return
	}

	filepath := flag.Arg(0)

	// handle --verify flag first
	if *verify {

		err := verifyFile(filepath)

		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
	// handle --batch second

	if *batch {
		err := batchHash(filepath)
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
>>>>>>> 6662348 (feat: add --batch flag for directory hashing)

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
	info, err := os.Stat(filepathArg)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Error: file not found")
		} else {
			fmt.Println("Error", err)
		}
		return
	}

	if info.IsDir() {
		fmt.Println("Error: that's a directory not a file. Use --batch to process directories.")
		return
	}

	fmt.Println("Name:", info.Name())
	fmt.Println("Size:", info.Size())
	fmt.Println("Modified:", info.ModTime().Format("2006-01-02 15:04:05"))

	hash, err := hashFile(filepathArg)
	if err != nil {
		fmt.Println("Error hashing file:", err)
		return
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
		fmt.Println("Error saving record:", err)
		return
	}
	fmt.Println("Record saved to record.json")
}
func verifyFile(filepath string) error {

<<<<<<< HEAD
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
=======
	//read the saved file

	data, err := os.ReadFile("record.json")

	if err != nil {

		return fmt.Errorf("cannot read record.json: %v", err)
	}
	//convert json back to struct

	var record FileRecord

	err = json.Unmarshal(data, &record)
	if err != nil {

		return fmt.Errorf("cannot parse record.json: %v", err)
	}
	//get current file's hash

	currentHash, err := hashFile(filepath)

	if err != nil {

		return fmt.Errorf("cannot hash file: %v", err)
	}
	if currentHash == record.SHA256 {
		fmt.Println("Recorded hash:", record.SHA256)
		fmt.Println("Current hash: ", currentHash)
		fmt.Println("VERIFIED: File has not been tampered with.")
	} else {
		fmt.Println("Recorded hash:", record.SHA256)
		fmt.Println("Current hash: ", currentHash)
>>>>>>> 6662348 (feat: add --batch flag for directory hashing)
		fmt.Println("ALERT: File has been modified since recording.")
	}
	return nil
}
<<<<<<< HEAD

=======
>>>>>>> 6662348 (feat: add --batch flag for directory hashing)
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

	jsonBytes, err := json.MarshalIndent(records, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile("batch_record.json", jsonBytes, 0o644)
	if err != nil {
		return err
	}

	fmt.Printf("Processed %d files. Record saved to batch_record.json\n", len(records))
	return nil
}
func batchHash(dirPath string) error {

	//create storage for multiple FileRecord structs
	//use a slice since its like a dynamic array that will grow

	//step 1: empty slice waiting to be filled

	var records []FileRecord

	//step 2: walk the directory

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		// step 3a: handle access errors
		if err != nil {

			return err
		}

		// step 3b: skip directories

		if info.IsDir() {

			return nil
		}

		// step 3c: Hash the file

		hash, err := hashFile(path)

		if err != nil {

			return err
		}
		// step 3d: build the record

		record := FileRecord{

			Filename:   info.Name(),
			Size:       info.Size(),
			Modified:   info.ModTime().Format("2006-01-02 15:04:05"),
			SHA256:     hash,
			RecordedAt: time.Now().Format("2006-01-02 15:04:05"),
			// step 3e add to slice
		}
		records = append(records, record)
		return nil

	})
	// step 4 handle any error from walk itself

	if err != nil {
		return err
	}
	//step 5: Marshal the slice to JSON
	// convert the records slice that contains all file records to JSON format
	bytes, err := json.MarshalIndent(records, "", "\t")

	if err != nil {
		return err
	}

	err = os.WriteFile("batch_record.json", bytes, 0644)

	if err != nil {

		return err
	}

	fmt.Printf("Processed %d files. Record saved to batch_record.json\n", len(records))
	return nil
}
