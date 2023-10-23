package main

import (
	"concurrency-pipeline/lib"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	log.Println("start")

	// pipeline 1 : read all file
	fileContent := lib.ReadFilesMysqlf()

	// pipeline 2 : compress all file
	chanFileCompressed1 := lib.Compress(fileContent)
	chanFileCompressed2 := lib.Compress(fileContent)
	chanFileCompressed3 := lib.Compress(fileContent)
	chanFileCompressed4 := lib.Compress(fileContent)
	chanFileCompressed := lib.MergeManyChannelMyself(chanFileCompressed1, chanFileCompressed2, chanFileCompressed3, chanFileCompressed4)

	// pipeline 3 : write compressed data to file
	chanFileWrited1 := lib.Write(chanFileCompressed)
	chanFileWrited2 := lib.Write(chanFileCompressed)
	chanFileWrited3 := lib.Write(chanFileCompressed)
	chanFileWrited4 := lib.Write(chanFileCompressed)
	chanFileWrited := lib.MergeManyChannelMyself(chanFileWrited1, chanFileWrited2, chanFileWrited3, chanFileWrited4)

	// pipeline 4 : search md5 sum
	chanFileSum1 := lib.GetSumMyself(chanFileWrited)
	chanFileSum2 := lib.GetSumMyself(chanFileWrited)
	chanFileSum3 := lib.GetSumMyself(chanFileWrited)
	chanFileSum := lib.MergeManyChannelMyself(chanFileSum1, chanFileSum2, chanFileSum3)

	// pipeline 5 : rename file
	chanFileRename1 := lib.RenameMyself(chanFileSum)
	chanFileRename2 := lib.RenameMyself(chanFileSum)
	chanFileRename3 := lib.RenameMyself(chanFileSum)
	chanFileRename := lib.MergeManyChannelMyself(chanFileRename1, chanFileRename2, chanFileRename3)

	// pipeline 6 : output

	counterTotal := 0
	counterRenamed := 0

	for fileInfo := range chanFileRename {
		if fileInfo.IsRename {
			counterRenamed++
		}
		counterTotal++
	}

	log.Printf("%d/%d file renamed\n", counterRenamed, counterTotal)

	duration := time.Since(start)
	fmt.Printf("%f seconds running\n", duration.Seconds())
}
