// file program ini adalah refactor dari file "2_find_md5_sum_of_file_then_rename_it"
// dimana menerapkan pattern pipeline dalam concurrency programming

package main

import (
	"log"
	"time"

	"concurrency-pipeline/lib"
)

func main() {
	log.Println("start")
	start := time.Now()

	// pipeline 1 : loop all files and read it
	chanFileContent := lib.ReadFiles()

	// pipeline 3 : calculate md5sum
	chanFileSum1 := lib.GetSum(chanFileContent)
	chanFileSum2 := lib.GetSum(chanFileContent)
	chanFileSum3 := lib.GetSum(chanFileContent)
	chanFileSum := lib.MergeChanFileInfo(chanFileSum1, chanFileSum2, chanFileSum3)

	// pipeline 4 : rename file
	chanRename1 := lib.Rename(chanFileSum)
	chanRename2 := lib.Rename(chanFileSum)
	chanRename3 := lib.Rename(chanFileSum)
	chanRename4 := lib.Rename(chanFileSum)
	chanRename := lib.MergeChanFileInfo(chanRename1, chanRename2, chanRename3, chanRename4)

	// pipeline 5 : output
	counterRenamed := 0
	counterTotal := 0

	for fileInfo := range chanRename {
		if fileInfo.IsRename {
			counterRenamed++
		}

		counterTotal++
	}

	log.Printf("%d/%d file renamed", counterRenamed, counterTotal)

	duration := time.Since(start)
	log.Printf("done in %.3f seconds", duration.Seconds())
}
