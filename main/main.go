package main

import "github.com/Ben-harder/dedupe"

func main() {
	source := "source.txt"
	dedupe.Concurrently(source, "out1.txt", 5)
	dedupe.RunAndWrite(source, "out2.txt")
}
