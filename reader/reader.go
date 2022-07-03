package reader

import (
	"bufio"
	"log"
	"os"
)

type reader struct {
	s         *bufio.Scanner
	f         *os.File
	offset    int64 // The offset to start reading from the file
	stop      int64 // The number of bytes to read before stopping
	bytesRead int64
	done      bool
}

// newReader opens the source file, seeks to offset bytes, and returns a scanner.
// The reader will cease reading once stop bytes have been read.
// Set stop to -1 to indicate EOF.
func New(path string, offset int64, stop int64) *reader {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(offset, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Initalizing reader to read %v bytes from offset %v\n", stop, offset)
	return &reader{
		s:      bufio.NewScanner(file),
		f:      file,
		offset: offset,
		stop:   stop,
	}
}

// Advance returns true if a line was able to be read and false if we are done or an error occurred
func (r *reader) Advance() bool {
	if r.bytesRead == 0 {
		log.Println("Starting read...")
	}
	done := r.done || !r.s.Scan()
	if done {
		log.Printf("Finished reading. Read %v bytes from offset %v", r.bytesRead, r.offset)

		// Check if read finished because of an error
		if err := r.s.Err(); err != nil {
			log.Fatal(err)
		}
		r.f.Close()
		return false
	}
	return true
}

// CurrLine returns the latest read line and increments the number of read bytes.
// Assumes the file scanner is split along newlines.
func (r *reader) CurrLine() string {
	r.bytesRead = r.bytesRead + int64(len(r.s.Text())) + 1 // + 1 is to account for newlines
	if r.stop != -1 && r.bytesRead >= r.stop {
		r.done = true
	}
	return r.s.Text()
}
