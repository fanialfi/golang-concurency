package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileInfoMysqlf struct {
	Err         error
	Index       int
	FileName    string
	WorkerIndex int // sebagai penanda worker mana yang akan melakukan operasi pembuatan file
}

func GenerateFileMyself() {
	os.RemoveAll(tempPath)
	os.MkdirAll(tempPath, os.ModePerm)

	// pipeline 1 : jalankan goruntine untuk job distribution
	chanFileIndex := generateFileIndexsMysqlf()

	// pipeline 2 : dispatch worker untuk menjalankan job
	chanFileResult1 := createFileMyself(chanFileIndex, 1)
	chanFileResult2 := createFileMyself(chanFileIndex, 2)
	chanFileResult3 := createFileMyself(chanFileIndex, 3)
	chanFileResult := mergeChanFile(chanFileResult1, chanFileResult2, chanFileResult3)

	// pipeline 3 : track and print out
	counterTotal := 0
	counterSucess := 0
	for fileResult := range chanFileResult {
		if fileResult.Err != nil {
			log.Printf("error creating file %s. stack trace %s\n", fileResult.FileName, fileResult.Err)
		} else {
			counterSucess++
		}
		counterTotal++
	}
	log.Printf("%d/%d of total files created", counterSucess, counterTotal)
}

func mergeChanFile(manyChanIn ...<-chan FileInfoMysqlf) <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)
	wg := new(sync.WaitGroup)

	wg.Add(len(manyChanIn))
	for _, eachChan := range manyChanIn {

		go func(eachChan <-chan FileInfoMysqlf) {
			defer wg.Done()

			for eachChanData := range eachChan {
				data <- eachChanData
			}
		}(eachChan)
	}

	go func() {
		wg.Wait()
		close(data)
	}()

	return data
}

func createFileMyself(chanIn <-chan FileInfoMysqlf, workerIndex int) <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	go func(workerIndex int) {
		defer close(data)

		for fileInfo := range chanIn {
			filePath := filepath.Join(tempPath, fileInfo.FileName)
			content := randomString(contentLength)

			err := os.WriteFile(filePath, []byte(content), os.ModePerm)
			// log.Printf("worker %d working on %s file generation\n", workerIndex, fileInfo.FileName)

			data <- FileInfoMysqlf{
				Err:         err,
				FileName:    fileInfo.FileName,
				WorkerIndex: workerIndex,
			}
		}
	}(workerIndex)

	return data
}

func generateFileIndexsMysqlf() <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	go func() {
		defer close(data)

		for i := 0; i < totalFile; i++ {
			data <- FileInfoMysqlf{
				Index:    i,
				FileName: fmt.Sprintf("file-%d.txt", i),
			}
		}
	}()

	return data
}
