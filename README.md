# Dedupe

Example usage:

```
func main() {
	source := "source.txt"
	dedupe.Concurrently(source, "out1.txt", 5)
	dedupe.RunAndWrite(source, "out2.txt")
}
```