// file program ini adalah refactor dari file "2_find_md5_sum_of_file_then_rename_it"
// dimana menerapkan pattern pipeline dalam concurrency programming
//
// dimana logic nya akan dipecah menjadi 3 bagian, dan seluruhnya di eksekusi secara konkuren
// - proses baca file
// - proses perhitungan md5 hash sum
// - proses rename file

package lib

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// FileInfo digunakan untuk metadata tiap file
// metadata ini diperlukan untuk mempermudah tracking file
type FileInfo struct {
	FilePath string // file location
	Content  []byte // file content
	Sum      string // md5 sum of content
	IsRename bool   // indicate whether the particular file is renamed already or not
}

// ReadFiles digunakan untuk pembacaan semua file, mengembalikan chan bertipe FileInfo
func ReadFiles() <-chan FileInfo {
	chanOut := make(chan FileInfo)

	// proses ini berjalan secara asynchronous dan concurrency
	go func() {
		// filepath.Walk akan mengiterasi seluruh content yang ada di folder tempPath
		err := filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
			// jika disini error, langsung return
			if err != nil {
				return err
			}

			// jika yang diiterasi sebuah sub directory, langsung return
			if info.IsDir() {
				return nil
			}

			// baca isi file
			buf, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// simpan / kirim informasi mengenai content (data file) dan path file kedalam channel
			chanOut <- FileInfo{
				FilePath: path,
				Content:  buf,
			}

			return nil
		})

		// jika function pada parameter kedua mengembalikan error
		if err != nil {
			log.Println("ERROR :", err.Error())
		}

		// jika semua file yang ada di folder tempPath sudah diiterasi semua, channel chanOut akan diclose
		close(chanOut)
	}()

	return chanOut
}

// GetSum digunakan untuk perhitungan md5 hash
// function ini juga biasa disebut dengan fan-out function
// fungsi fan-out digunakan untuk pendistribusian job ke banyak worker.
// multiple function bisa membaca dari channel yang sama sampai channel yang dibaca closed, ini disebut fan-out
// parameter bertipe channel pada function GetSum merupakan media untuk distribusi job
// sedangkan setiap pemanggilan function GetSum sendiri merepresentasikan sebagai 1 worker
func GetSum(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		// setiap ada penerimaan data baru dari chanIn akan di listen, dan kemudian dilanjutkan ke proses kalkulasi md5 sum
		for fileInfo := range chanIn {
			// hasil dari md5 sum ditambahkan ke data FileInfo field sum
			fileInfo.Sum = fmt.Sprintf("%x", md5.Sum(fileInfo.Content))

			// kemudian dikirim lagi ke chanOut, yang mana channel ini merupakan nilai balik dari function GetSum
			chanOut <- fileInfo
		}

		// ketika chanIn closed maka diasumsikan semua data sudah dikirim, dan chanOut juga akan di close
		close(chanOut)
	}()

	return chanOut
}

// MergeChanFileInfo menggabungkan banyak channel ke satu channel saja
// dimana channel ini (chanOut didalam function MergeChanFileInfo) juga akan otomatis diclose ketika channel inputan, channel pada argument function adalah closed
// function seperti ini biasa disebut fan-in function
func MergeChanFileInfo(chanInMany ...<-chan FileInfo) <-chan FileInfo {
	wg := new(sync.WaitGroup)
	chanOut := make(chan FileInfo)

	wg.Add(len(chanInMany))
	for _, eachChan := range chanInMany {

		go func(eachChan <-chan FileInfo) {
			for eachChanData := range eachChan {
				chanOut <- eachChanData
			}
			wg.Done()
		}(eachChan)

	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

// Rename secara garis besar penulisan function hampir sama dengan GetSum
// hanya saja fungsinya untuk rename file
func Rename(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		for fileInfo := range chanIn {
			newPath := filepath.Join(tempPath, fmt.Sprintf("file-%s.txt.gz", fileInfo.Sum))

			err := os.Rename(fileInfo.FilePath, newPath)

			// lakukan pengecekan apakah variabel err nilainya nil
			// hasil dari pengecekain ini disimpan ke field IsRename
			fileInfo.IsRename = (err == nil)
			chanOut <- fileInfo
		}

		close(chanOut)
	}()

	return chanOut
}
