package lib

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

type FileInfo struct {
	Err         error
	Index       int
	FileName    string
	WorkerIndex int
}

func GenerateFiles() {
	os.RemoveAll(tempPath)
	os.MkdirAll(tempPath, os.ModePerm)

	// pipeline 1 : job distribution
	chanFileIndex := generateFilesIndex()

	// pipeline 2 : the main logic (creating file)
	chanFilesWorker := 100
	chanFileResult := createFile(chanFileIndex, chanFilesWorker)

	// pipeline 3 : track and print input
	counterTotal := 0
	counterSuccess := 0
	for fileResult := range chanFileResult {
		if fileResult.Err != nil {
			log.Printf("error creating file %s. stack trace %s\n", fileResult.FileName, fileResult.Err.Error())
		} else {
			counterSuccess++
		}
		counterTotal++
	}

	log.Printf("%d/%d of total files created\n", counterSuccess, counterTotal)
}

func createFile(chanIn <-chan FileInfo, numberOfWorker int) <-chan FileInfo {
	data := make(chan FileInfo)
	wg := new(sync.WaitGroup)

	wg.Add(numberOfWorker)
	go func() {

		for workerIndex := 0; workerIndex < numberOfWorker; workerIndex++ {
			go func(workerIndex int) {
				defer wg.Done()

				for job := range chanIn {
					filepath := filepath.Join(tempPath, job.FileName)
					content := randomString(contentLength)

					err := os.WriteFile(filepath, []byte(content), os.ModePerm)

					log.Printf("worker %d working on %s file generation\n", workerIndex, job.FileName)

					data <- FileInfo{
						FileName:    job.FileName,
						WorkerIndex: workerIndex,
						Err:         err,
					}
				}

			}(workerIndex)
		}

	}()

	go func() {
		wg.Wait()
		close(data)
	}()

	return data
}

func generateFilesIndex() <-chan FileInfo {
	data := make(chan FileInfo)

	go func() {
		defer close(data)

		for i := 0; i < totalFile; i++ {
			data <- FileInfo{
				Index:    i,
				FileName: fmt.Sprintf("file-%d.txt", i),
			}
		}
	}()

	return data
}

func randomString(length int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRTSUVWXYZ")

	b := make([]rune, length)
	for index := range b {
		b[index] = letter[rand.Intn(len(letter))]
	}

	return string(b)
}
