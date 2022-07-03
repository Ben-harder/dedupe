package writer

import (
	"bufio"
	"log"
	"os"
	"sync"
)

type Writer struct {
	w  *bufio.Writer
	f  *os.File
	mu sync.Mutex
}

func New(path string) *Writer {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	return &Writer{
		w: bufio.NewWriter(file),
		f: file,
	}
}

func (w *Writer) WriteLine(line string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := w.w.WriteString(line + "\n")
	if err != nil {
		log.Fatal(err)
	}
	w.w.Flush()
}

func (w *Writer) Close() {
	err := w.f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
