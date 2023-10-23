package main

import (
	"contex-cancelation-pipeline/lib"
	"log"
	"time"
)

func main() {
	log.Println("start")
	start := time.Now()

	lib.GenerateFiles()

	duration := time.Since(start)
	log.Printf("done in %f seconds\n", duration.Seconds())
}
