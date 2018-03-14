package main

import (
	"rsync_go/file_sync_client"
)

func main() {
	//TODO: Directories should come from command-line variables
	srcDirectory := "files/main"
	destDirectory := "files/copy"

	client, err := sync.NewFileSyncClient(srcDirectory, destDirectory)
	if err != nil {
		panic(err)
	}

	client.Sync()
}
