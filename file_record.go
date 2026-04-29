package main

type FileRecord struct {
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	Modified   string `json:"modified"`
	SHA256     string `json:"sha256"`
	RecordedAt string `json:"recorded_at"`
}
