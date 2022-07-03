package dedupe

import (
	"log"
	"os"

	"github.com/Ben-harder/dedupe/reader"
	"github.com/Ben-harder/dedupe/set"
	"github.com/Ben-harder/dedupe/writer"
)

// RunAndWrite will dedupe the source file while simultaneously writing distinct lines to the output file.
func RunAndWrite(sourceFile, outputFile string) {
	var (
		d = newDeduper(sourceFile)
		r = reader.New(d.src, 0, -1)
		w = writer.New(outputFile)
	)

	for r.Advance() {
		line := r.CurrLine()
		if d.isNewLine(line) {
			w.WriteLine(line)
		}
	}
	w.Close()
}

// Concurrently will divide the source file into even portions based on the number of workers
// and deduplicate them concurrently. The results are combined and written to the output file.
func Concurrently(sourceFile, outputFile string, workers int) {
	setChan := make(chan *set.Set)

	// Start workers
	dedupeConcurrently(sourceFile, workers, setChan)

	// Combine work sequentially as it is completed. Doing this concurrently could be an optimization opportunity if desired.
	combined := set.CombineFromChan(workers, setChan)

	// Output to file
	dumpToFile(combined, outputFile)
}

// Starts workers to build deduplicated sets concurrently.
// When finished, the set is sent to setChan.
func dedupeConcurrently(sourceFile string, workers int, setChan chan *set.Set) {
	offsets := divideFile(sourceFile, workers)
	log.Println("Divided file for processing:", offsets)

	for i := 0; i < workers; i++ {
		go func(i int) {
			d := newDeduper(sourceFile)
			if i == workers-1 {
				d.buildDedupedSet(offsets[i], -1)
			} else {
				d.buildDedupedSet(offsets[i], offsets[1])
			}
			setChan <- d.dedupedSet
		}(i)
	}
}

type deduper struct {
	src        string
	dedupedSet *set.Set
}

func newDeduper(src string) *deduper {
	return &deduper{
		src:        src,
		dedupedSet: set.New(),
	}
}

// BuildDedupedSet builds a set of unique lines from the source file.
// Starts processing the source file from byte "from" until "stop" bytes
// have been processed.
// Use -1 for "to" to indicate EOF.
func (d *deduper) buildDedupedSet(from, stop int64) {
	var (
		r = reader.New(d.src, from, stop)
	)
	for r.Advance() {
		line := r.CurrLine()
		d.process(line)
	}
}

// process inserts a line into the deduped set
func (d *deduper) process(line string) {
	d.dedupedSet.Insert(line)
}

// isNewLine returns true if we have not come across currLine yet and inserts it into the deduped set.
func (d *deduper) isNewLine(currLine string) bool {
	if !d.dedupedSet.Contains(currLine) {
		d.process(currLine)
		return true
	}
	return false
}

// divideFile takes a file, divides its size, and returns the offsets.
// Ex: 1000 byte file with n = 2 returns [0, 500]
func divideFile(file string, n int) []int64 {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	stats, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	offsets := make([]int64, 0)
	for i := 0; i < n; i++ {
		offsets = append(offsets, 0+int64(i)*(stats.Size()/int64(n)))
	}
	return offsets
}

// dumpToFile writes out the keys in dedupedSet to file.
func dumpToFile(dedupedSet *set.Set, file string) {
	w := writer.New(file)
	for _, k := range dedupedSet.Keys() {
		w.WriteLine(k)
	}
	w.Close()
}
