package main

import (
	"concurrency-pipeline/lib"
	"log"
	"time"
)

func main() {
	log.Println("start")
	start := time.Now()

	lib.GenerateFile()

	duration := time.Since(start)
	log.Printf("done in %.3f second\n", duration.Seconds())
}
