package chunker

import (
	"bufio"
	"crypto/md5"
	"io/ioutil"
	"os"
)

func computeSumHashDigest(data []byte) int64 {
	var sum int64

	for i := range data {
		sum += int64(data[i])
	}

	return sum
}

type Chunk struct {
	Data            []byte
	WeakSignature   int64
	StrongSignature [16]byte
}

func NewChunk(data []byte) *Chunk {
	return &Chunk{Data: data}
}

func (chunk *Chunk) ComputeWeakSignature() {
	chunk.WeakSignature = computeSumHashDigest(chunk.Data)
}

func (chunk *Chunk) ComputeStrongSignature() {
	chunk.StrongSignature = md5.Sum(chunk.Data)
}

const (
	BLOCK_BOUNDARY = iota
	BYTE_BOUNDARY
)

type chunking interface {
	Chunk() []*Chunk
}

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
	bytes := getFileBytes(chunker.fd)

	chunkCnt := len(bytes) - chunker.blockSize

	chunks := make([]*Chunk, 0)
	for index := 0; index < chunkCnt; index++ {
		endOffset := index + chunker.blockSize
		chunkBytes := bytes[index:endOffset]

		chunks = append(chunks, NewChunk(chunkBytes))
	}

	return chunks
}

func (chunker *Chunker) ChunksAtBlockBoundary() []*Chunk {
	bytes := getFileBytes(chunker.fd)

	chunkCnt := len(bytes) / chunker.blockSize
	if len(bytes)%chunker.blockSize > 0 {
		chunkCnt += 1
	}

	chunks := make([]*Chunk, 0)
	for index := 0; index < chunkCnt-1; index++ {
		startOffset := index * chunker.blockSize
		endOffset := (index + 1) * chunker.blockSize
		chunkBytes := bytes[startOffset:endOffset]

		chunks = append(chunks, NewChunk(chunkBytes))
	}

	lastStartOffset := (chunkCnt - 1) * chunker.blockSize
	chunks = append(chunks, NewChunk(bytes[lastStartOffset:]))

	return chunks
}

func getFileBytes(fd *os.File) []byte {
	reader := bufio.NewReader(fd)

	bytes, err := ioutil.ReadAll(reader)
	check(err)

	return bytes
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
