package chunker

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

type Chunk struct {
	data            []byte
	weakSignature   []byte
	strongSignature []byte
}

func NewChunk(data []byte) *Chunk {
	return &Chunk{data: data}
}

const (
	BLOCK_BOUNDARY = iota
	BYTE_BOUNDARY
)

type Chunker struct {
	fd        *os.File
	blockSize int // in bytes
}

func NewChunker(fd *os.File, blockSize int) *Chunker {
	return &Chunker{fd: fd, blockSize: blockSize}
}

func (chunker *Chunker) Chunks(mode int) []*Chunk {
	if mode == BLOCK_BOUNDARY {
		return chunker.ChunksAtBlockBoundary()
	} else if mode == BYTE_BOUNDARY {
		return chunker.ChunksAtByteBoundary()
	} else {
		// unknown
		return make([]*Chunk, 0)
	}
}

func (chunker *Chunker) ChunksAtByteBoundary() []*Chunk {
	fmt.Println("Chunking at byte boundaries is not implemented yet.")
	return make([]*Chunk, 0)
}

func (chunker *Chunker) ChunksAtBlockBoundary() []*Chunk {
	reader := bufio.NewReader(chunker.fd)

	bytes, err := ioutil.ReadAll(reader)
	check(err)

	chunkCnt := len(bytes) / chunker.blockSize
	if len(bytes)%chunker.blockSize > 0 {
		chunkCnt += 1
	}

	chunks := make([]*Chunk, 0)
	for index := 0; index < chunkCnt-1; index++ {
		startOffset := index * chunker.blockSize
		endOffset := (index + 1) * chunker.blockSize
		chunks = append(chunks, NewChunk(bytes[startOffset:endOffset]))
	}

	lastStartOffset := (chunkCnt - 1) * chunker.blockSize
	chunks = append(chunks, NewChunk(bytes[lastStartOffset:]))

	return chunks
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
