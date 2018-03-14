package main

import (
	"rsync_go/file_sync_client"
)

func main() {
	srcDirectory := "files/main"
	destDirectory := "files/copy"

	client, _ := sync.NewFileSyncClient(srcDirectory, destDirectory)
	client.Sync()
}
