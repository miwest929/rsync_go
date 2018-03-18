package sync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"rsync_go/file_sync"
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

func (client *FileSyncClient) Sync() {
	// Get all files in source directory
	// Get all files in destination directory
	// If file exists in source and not destination then copy entire file over
	// If file exists in both then do following:
	//   Split source version into chunks. Split at block boundary
	srcFiles := getDirectoryFiles(client.srcDirectory)
	destFiles := getDirectoryFiles(client.destDirectory)

	defer closeFiles(srcFiles)
	defer closeFiles(destFiles)

	destFilesSet := make(map[string]*os.File)
	for _, f := range destFiles {
		baseName := filepath.Base(f.Name())
		destFilesSet[baseName] = f
	}

	for _, srcFile := range srcFiles {
		baseName := filepath.Base(srcFile.Name())
		destFile, ok := destFilesSet[baseName]
		if ok {
			fmt.Printf("Syncing '%s' with '%s'\n", srcFile.Name(), destFile.Name())
			client.performSync(srcFile, destFile)
		} else {
			// file doesn't exist in destination directory. Must copy it over.
			fmt.Printf("INFO: '%s' doesn't exist in source directory.", srcFile.Name())
		}
	}
}

func (client *FileSyncClient) performSync(srcFd *os.File, destFd *os.File) {
	fileSync := file_sync.NewFileSync(srcFd, destFd)
	fileSync.Sync()
	//srcChunk.ComputeWeakSignature()
	//srcChunk.ComputeStrongSignature()

	// The slow hash for the destination file is computing in a rolling fashion
	// This is called a "rolling hash"
	//destChunk.ComputeWeakSignature()
}

func closeFiles(files []*os.File) {
	for _, file := range files {
		file.Close()
	}
}

func getDirectoryFiles(directory string) []*os.File {
	filenames := []string{}

	//TODO: Consider using a faster way to crawl a directory than filepath.Walk
	//TODO: Verify if Walk does a recursive crawl. Make the recursive behavior optional
	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		// Only non-directory files will be synched
		if !f.IsDir() {
			filenames = append(filenames, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	files := []*os.File{}
	for _, filename := range filenames {
		fd, err := os.Open(filename)

		if err != nil {
			panic(err)
		}

		files = append(files, fd)
	}

	return files
}

func isDirectory(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}
