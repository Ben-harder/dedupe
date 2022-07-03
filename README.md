# Dedupe

Go package to deduplicate lines of text files.

Example usage:

```
func main() {
	source := "source.txt"

	// Concurrently deduplicates lines in the source file using 5 Goroutines and outputs them to out1.txt
	dedupe.Concurrently(source, "out1.txt", 5)

	// In one go, deduplicates lines in the source file while writing distinct lines to out2.txt 
	dedupe.RunAndWrite(source, "out2.txt")
}
```