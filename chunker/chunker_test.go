package chunker

import (
	"os"
	"rsync_go/chunker"
	"testing"
)

const (
	BLOCK_SIZE = 8
)

func TestBoundaryChunker(t *testing.T) {
	fd := openTestFile("../files/test/file.txt")

	c := chunker.NewChunker(fd, BLOCK_SIZE)
	chunks := c.Chunks(chunker.BLOCK_BOUNDARY)

	if len(chunks) != 3 {
		t.Fatalf("Expect %d chunks. Got %d.", 3, len(chunks))
	}
}

func TestByteChunker(t *testing.T) {
	fd := openTestFile("../files/test/file.txt")

	c := chunker.NewChunker(fd, BLOCK_SIZE)
	chunks := c.Chunks(chunker.BYTE_BOUNDARY)

	if len(chunks) != 9 {
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
