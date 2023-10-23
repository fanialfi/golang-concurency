package lib

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

type FileInfoWithContext struct {
	Err         error
	Index       int
	FileName    string
	WorkerIndex int
}

func GenerateFilesWithContext(ctx context.Context) {
	os.RemoveAll(tempPath)
	os.MkdirAll(tempPath, os.ModePerm)

	done := make(chan int) // sebagai indikator bahwa proses pipeline sudah selesai secara keseluruhan

	go func() {

		// pipeline 1 : job distribution
		chanFileIndex := generateFilesIndexWithContext(ctx)

		// pipeline 2 : the main logic (creating file)
		chanFilesWorker := 100
		chanFileResult := createFileWithContext(ctx, chanFileIndex, chanFilesWorker)

		// pipeline 3 : track and print input
		counterSuccess := 0
		for fileResult := range chanFileResult {
			if fileResult.Err != nil {
				log.Printf("error creating file %s. stack trace %s\n", fileResult.FileName, fileResult.Err.Error())
			} else {
				counterSuccess++
			}
		}

		// notify bahwa proses sudah selesai
		// mengirim informasi jumplah file yang berhasil di generate
		done <- counterSuccess
	}()

	select {

	// ketika channel Done milik context ini menerima data, berarti context telah di cancle secara paksa
	// cancle nya bisa karena memang context sudah melebihi timeout yang sudah ditentukan
	// atau dicancle secara explisit lewat pemanggilan callback contect.CancleFunc
	case <-ctx.Done():
		// untuk mengetahui alasan cancle, bisa dengan mengakses method Err milik context (ctx.Err)
		log.Printf("generation process stopped %s", ctx.Err().Error())

		// case kedua ini akan terpenuhi ketika proses pipeline sudah selesai secara keseluruhan
	case counterSuccess := <-done:
		log.Printf("%d/%d of total file created", counterSuccess, totalFile)
	}
}

// cancellation juga perlu diterapkan di createFileWithContext
// jika tidak, proses pembuatan file akan tetap berjalan sesuai dengan jumplah job yang dikirim meskipun sudah dicancel secara paksa
func createFileWithContext(ctx context.Context, chanIn <-chan FileInfoWithContext, numberOfWorker int) <-chan FileInfoWithContext {
	data := make(chan FileInfoWithContext)
	wg := new(sync.WaitGroup)

	wg.Add(numberOfWorker)
	go func() {

		for workerIndex := 0; workerIndex < numberOfWorker; workerIndex++ {
			go func(workerIndex int) {
				defer wg.Done()

				for job := range chanIn {
					select {
					case <-ctx.Done():
						break
					default:
						filepath := filepath.Join(tempPath, job.FileName)
						content := randomStringWithContext(contentLength)

						err := os.WriteFile(filepath, []byte(content), os.ModePerm)

						log.Printf("worker %d working on %s file generation\n", workerIndex, job.FileName)

						data <- FileInfoWithContext{
							FileName:    job.FileName,
							WorkerIndex: workerIndex,
							Err:         err,
						}
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

// meskipun GenerateFilesIndexWithContext otomatis di stop ketika cancelled
// proses didalamnya akan tetap berjalan jika tidak di handle dengan baik cancellation nyax
func generateFilesIndexWithContext(ctx context.Context) <-chan FileInfoWithContext {
	data := make(chan FileInfoWithContext)

	go func() {
		defer close(data)

		for i := 0; i < totalFile; i++ {
			select {
			// jika ada notif cancel secara paksa, maka case terpenuhi, dan berulangan di break
			case <-ctx.Done():
				break
			default:
				data <- FileInfoWithContext{
					Index:    i,
					FileName: fmt.Sprintf("file-%d.txt", i),
				}
			}
		}
	}()

	return data
}

func randomStringWithContext(length int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRTSUVWXYZ")

	b := make([]rune, length)
	for index := range b {
		b[index] = letter[rand.Intn(len(letter))]
	}

	return string(b)
}
