package main

import (
	"log"
	"runtime"
	"simplified-fan-out-fan-in-pipeline/lib"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	log.Println("start")
	start := time.Now()

	lib.GenerateFileMyself()

	duration := time.Since(start)
	log.Printf("done in %f seconds\n", duration.Seconds())
}
