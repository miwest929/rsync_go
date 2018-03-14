package sync

import (
	"errors"
	"os"
	"rsync_go/chunker"
)

type FileSyncClient struct {
	srcDirectory  string
	destDirectory string
}

func NewFileSyncClient(srcDirectory string, destDirectory string) (*FileSyncClient, error) {
	if !isDirectory(srcDirectory) {
		return nil, errors.New("src directory is not actually a directory")
	}

	if !isDirectory(destDirectory) {
		return nil, errors.New("destination directory is not actually a directory")
	}

	return &FileSyncClient{srcDirectory: srcDirectory, destDirectory: destDirectory}, nil
}

func isDirectory(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}

func (client *FileSyncClient) Sync() {
	// Get all files in source directory
	// Get all files in destination directory
	// If file exists in source and not destination then copy entire file over
	// If file exists in both then do following:
	//   Split source version into chunks. Split at block boundary
}
