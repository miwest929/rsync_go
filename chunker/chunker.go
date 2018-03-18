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

type chunking interface {
	Chunk() []*Chunk
}

type BlockChunker struct {
	fd        *os.File
	blockSize int // in bytes
}

func NewBlockChunker(fd *os.File, blockSize int) *BlockChunker {
	return &BlockChunker{fd: fd, blockSize: blockSize}
}

func (chunker *BlockChunker) Chunks() []*Chunk {
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

type ByteChunker struct {
	fd        *os.File
	blockSize int // in bytes
}

func NewByteChunker(fd *os.File, blockSize int) *ByteChunker {
	return &ByteChunker{fd: fd, blockSize: blockSize}
}

func (chunker *ByteChunker) Chunks() []*Chunk {
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
