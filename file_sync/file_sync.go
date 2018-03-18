/*
  This implements a specific algorithm for file synchronization
*/
package file_sync

import (
	"fmt"
	"os"
	"rsync_go/chunker"
)

const (
	BLOCK_SIZE = 2000
)

type RollingHash interface {
	// Compute digest for the first block
	FirstHash() int64

	// Given the a digest for a block compute the digest for the successive block
	NextHash(int64, int) int64
}

type BatchSumHash struct {
	chunks []*chunker.Chunk
}

func (hash *BatchSumHash) FirstHash() int64 {
	chunk := hash.chunks[0]
	var sum int64

	for i := range chunk.Data {
		sum += int64(chunk.Data[i])
	}

	return sum
}

func (hash *BatchSumHash) NextHash(currDigest int64, chunkIndex int) int64 {
	if chunkIndex+1 >= len(hash.chunks) {
		return 0
	}

	chunk := hash.chunks[chunkIndex+1]
	var sum int64

	for i := range chunk.Data {
		sum += int64(chunk.Data[i])
	}

	return sum
}

type RollingSumHash struct {
	chunks []*chunker.Chunk
}

func (hash *RollingSumHash) FirstHash() int64 {
	chunk := hash.chunks[0]
	var sum int64

	for i := range chunk.Data {
		sum += int64(chunk.Data[i])
	}

	return sum
}

func (hash *RollingSumHash) NextHash(currDigest int64, chunkIndex int) int64 {
	if chunkIndex+1 >= len(hash.chunks) {
		return 0
	}

	currChunk := hash.chunks[chunkIndex]
	nextChunk := hash.chunks[chunkIndex+1]

	nextDigest := currDigest - int64(currChunk.Data[0])
	return nextDigest + int64(nextChunk.Data[len(nextChunk.Data)-1])
}

type FileSync struct {
	srcFile     *os.File
	destFile    *os.File
	blockChunks []*chunker.Chunk
	byteChunks  []*chunker.Chunk
}

func NewFileSync(srcFile *os.File, destFile *os.File) *FileSync {
	blockChunker := chunker.NewChunker(srcFile, BLOCK_SIZE)
	byteChunker := chunker.NewChunker(destFile, BLOCK_SIZE)

	return &FileSync{
		blockChunks: blockChunker.Chunks(chunker.BLOCK_BOUNDARY),
		byteChunks:  byteChunker.Chunks(chunker.BYTE_BOUNDARY),
		srcFile:     srcFile,
		destFile:    destFile,
	}
}

func (fileSync *FileSync) Sync() {
	// Here is the core of the file synchronization algorithm

	// Compute weak & strong hashes for srcFile chunks
	var blockWeakHash RollingHash
	blockWeakHash = &BatchSumHash{chunks: fileSync.blockChunks}
	currHash := blockWeakHash.FirstHash()
	lastBlockChunkIndex := len(fileSync.blockChunks) - 1
	for blockIndex, _ := range fileSync.blockChunks[0 : lastBlockChunkIndex-1] {
		currHash := blockWeakHash.NextHash(currHash, blockIndex)
		fmt.Printf("Chunk %d: %d\n", blockIndex, currHash)
	}

	// Compute weak hash for destFile chunks
	var weakHash RollingHash
	weakHash = &RollingSumHash{chunks: fileSync.byteChunks}
	currWeakHash := weakHash.FirstHash()
	lastChunkIndex := len(fileSync.byteChunks) - 1
	for blockIndex, _ := range fileSync.byteChunks[0 : lastChunkIndex-1] {
		currWeakHash = weakHash.NextHash(currWeakHash, blockIndex)

		// Search for weakHash in srcFile...if match then ..you know..
		fmt.Printf("Chunk %d: %d\n", blockIndex, currWeakHash)
	}
}
