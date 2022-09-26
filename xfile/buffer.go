package xfile

// https://google.golang.org/api 里面cv来的

import (
	"bytes"
	"io"
)

// Buffer buffers data from an io.Reader to support uploading media in
// retryable chunks. It should be created with NewBuffer.
type Buffer struct {
	media io.Reader

	chunk []byte // The current chunk which is pending upload.  The capacity is the chunk size.
	err   error  // Any error generated when populating chunk by reading media.

	// The absolute position of chunk in the underlying media.
	off int64
}

// NewBuffer initializes a Buffer.
func NewBuffer(media io.Reader, chunkSize int) *Buffer {
	return &Buffer{media: media, chunk: make([]byte, 0, chunkSize)}
}

// Chunk returns the current buffered chunk, the offset in the underlying media
// from which the chunk is drawn, and the size of the chunk.
// Successive calls to Chunk return the same chunk between calls to Next.
func (mb *Buffer) Chunk() (chunk io.Reader, off int64, size int, err error) {
	// There may already be data in chunk if Next has not been called since the previous call to Chunk.
	if mb.err == nil && len(mb.chunk) == 0 {
		mb.err = mb.loadChunk()
	}
	return bytes.NewReader(mb.chunk), mb.off, len(mb.chunk), mb.err
}

// loadChunk will read from media into chunk, up to the capacity of chunk.
func (mb *Buffer) loadChunk() error {
	bufSize := cap(mb.chunk)
	mb.chunk = mb.chunk[:bufSize]

	read := 0
	var err error
	for err == nil && read < bufSize {
		var n int
		n, err = mb.media.Read(mb.chunk[read:])
		read += n
	}
	mb.chunk = mb.chunk[:read]
	return err
}

// Next advances to the next chunk, which will be returned by the next call to Chunk.
// Calls to Next without a corresponding prior call to Chunk will have no effect.
func (mb *Buffer) Next() {
	mb.off += int64(len(mb.chunk))
	mb.chunk = mb.chunk[0:0]
}
