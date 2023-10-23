package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"simplified-fan-out-fan-in-pipeline/lib"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	runtime.GOMAXPROCS(2)
	log.Println("start")
	start := time.Now()

	lib.GenerateFileSequentinally()

	duration := time.Since(start)
	log.Printf("done in %f seconds\n", duration.Seconds())
}
