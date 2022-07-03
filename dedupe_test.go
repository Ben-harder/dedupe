package dedupe

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
)

const (
	mb         = 1000000
	sourceFile = "sample.txt"
	outputFile = "output.txt"
)

func TestConcurrently(t *testing.T) {
	Concurrently(sourceFile, outputFile, 5)
}

func TestSingleThreaded(t *testing.T) {
	RunAndWrite(sourceFile, outputFile)
}

func BenchmarkConcurrent10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Concurrently(sourceFile, outputFile, 10)
	}
}
func BenchmarkConcurrent5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Concurrently(sourceFile, outputFile, 5)
	}
}

func generateSourceFile(name string) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	letters := make([]string, 0)
	for r := 'a'; r <= 'z'; r++ {
		letters = append(letters, fmt.Sprintf("%c", r))
	}

	for i := 0; i < 5e7; i++ {
		_, err := writer.WriteString(letters[rand.Intn(26)] + "\n")
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	}
}
