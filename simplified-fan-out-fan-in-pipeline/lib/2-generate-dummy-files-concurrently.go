package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// FileInfo sebagai skema payload data ketika dikirim via channel dari goruntine jobs ke goruntine worker
type FileInfo struct {
	Err         error
	Index       int
	FileName    string
	WorkerIndex int // sebagai penanda worker mana yang akan melakukan operasi pembuatan file
}

func GenerateFileConcurency() {
	os.RemoveAll(tempPath)
	os.MkdirAll(tempPath, os.ModePerm)

	// pipeline 1 = dispatch goruntine for job distribution
	chanFileIndex := generateFileIndexes()

	// pipeline 2 = the main logic (creating file) / dispatch goruntine untuk start worker,
	// masing maisng worker bertugas untuk membuat file
	createFilesWorker := 3
	chanFileResult := createFile(chanFileIndex, createFilesWorker)

	// track and print output
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

// generateFileIndexes digunakan untuk membuat nama file dan untuk distribusi jobs dengan men dispatch goruntine
func generateFileIndexes() <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		for i := 0; i < totalFile; i++ {
			chanOut <- FileInfo{
				Index:    i,
				FileName: fmt.Sprintf("file-%d.txt", i),
			}
		}

		close(chanOut)
	}()

	return chanOut
}

// createFile merupakan function Fan-Out Fan-In karena menerima parameter channel sebelumnya
// createFile men dispatch goruntine worker dan men track untuk setiap output dari masing masing worker ke channel output
// createFile menghasilkan channel yang isinya adalah hasil dari masing masing operasi goruntine worker
func createFile(chanIn <-chan FileInfo, numberOfWorker int) <-chan FileInfo {
	// sebagai channel output Fan-In dari worker worker yang ada
	// secara langsung digunakan untuk return value dari function createFile
	// karena selain deklarasi channel dai WaitGroup, semua berjalan secara asynchronous menggunakan goruntine
	chanOut := make(chan FileInfo)

	// wait group untuk keperluan manajemen worker
	wg := new(sync.WaitGroup)

	// allocation N of worker
	wg.Add(numberOfWorker)

	go func() {

		// dispatch N worker
		for workerIndex := 0; workerIndex < numberOfWorker; workerIndex++ {

			go func(workerIndex int) {
				// listen to chanIn channel for incoming jobs
				for job := range chanIn {

					// do the job
					filepath := filepath.Join(tempPath, job.FileName)
					content := randomString(contentLength)

					err := os.WriteFile(filepath, []byte(content), os.ModePerm)

					// log.Printf("worker %d working on %s file generation\n", workerIndex, job.FileName)

					// construct the job's result, and send it to chanOut
					chanOut <- FileInfo{
						FileName:    job.FileName,
						WorkerIndex: workerIndex,
						Err:         err,
					}
				}

				// if chan is closed, and the remaning jobs are finished
				// only then we mark the worker as complete
				wg.Done()
			}(workerIndex)
		}

	}()

	// wait until chanIn closed and then all workers are done
	// because right after that - we need to close the chanOut channel
	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}
