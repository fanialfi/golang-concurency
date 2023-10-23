package main

import (
	"concurrency-pipeline/lib"
	"log"
	"time"
)

func main() {
	log.Println("starting")
	start := time.Now()

	lib.Proceed()

	duration := time.Since(start)
	log.Printf("done in %.3f seconds\n", duration.Seconds())
}
