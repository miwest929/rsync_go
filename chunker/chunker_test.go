package chunker

import (
	"os"
	"rsync_go/chunker"
	"testing"
)

func TestChunker(t *testing.T) {
	contents := []byte("ThisIsAByteArray")
	fd := openTestFile("../files/test/file.txt")

	numOfBlocks := 2
	blockSize := len(contents) / numOfBlocks
	c := chunker.NewChunker(fd, blockSize)
	chunks := c.Chunks(chunker.BLOCK_BOUNDARY)

	if len(chunks) != 3 {
		t.Fatalf("Expect %d chunks. Got %d.", 3, len(chunks))
	}
}

func openTestFile(filename string) *os.File {
	testFd, err := os.Open(filename)
	check(err)

	return testFd
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
