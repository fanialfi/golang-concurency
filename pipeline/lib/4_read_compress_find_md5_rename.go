package lib

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileInfoMysqlf struct {
	FilePath string // file location
	Content  []byte // file content
	Sum      string // md5 sum
	IsRename bool   //  mengindikasi bahwa file yang terkait sudah di rename
}

// ReadFilesMysqlf digunakan untuk membaca informasi tiap file yang diiterasi,
// metadata dari file yang diiterasi kemudian di kirim kedalam channel
// channel kemudian digunakan sebagai return value
func ReadFilesMysqlf() <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	go func() {
		defer close(data)

		err := filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// jika yang diiterasi adalah sebuah directory, langsung return
			if info.IsDir() {
				return nil
			}

			// baca file
			buffer, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// simpan informasi mengenai file yang sedang diiterasi kedalah channel
			data <- FileInfoMysqlf{
				FilePath: path,
				Content:  buffer,
			}

			return nil
		})

		if err != nil {
			log.Println("Error :", err.Error())
		}

	}()

	return data
}

// Compress digunakan untuk mengkompress data dari FileInfoMysqlf field Content
// kemudian hasil dari compress data di rewrite ke FileInfoMysqlf field Content
func Compress(chanIn <-chan FileInfoMysqlf) <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	go func() {
		defer close(data)

		for fileInfo := range chanIn {
			var bfr bytes.Buffer

			gzWriter, err := gzip.NewWriterLevel(&bfr, gzip.BestCompression)
			if err != nil {
				log.Println("Error", err.Error())
			}

			// compress data
			_, err = gzWriter.Write(fileInfo.Content)
			if err != nil {
				log.Println("Error", err.Error())
			}

			gzWriter.Close()

			// save data ke FileInfoMysqlf field Content
			fileInfo.Content = bfr.Bytes()
			// err = os.WriteFile(fileInfo.FilePath, bfr.Bytes(), os.ModePerm)
			// if err != nil {
			// 	log.Println(err.Error())
			// }

			data <- fileInfo
		}
	}()

	return data
}

func Write(chanIn <-chan FileInfoMysqlf) <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	go func() {
		defer close(data)

		for fileInfo := range chanIn {
			err := os.WriteFile(fileInfo.FilePath, fileInfo.Content, os.ModePerm)
			if err != nil {
				log.Println(err.Error())
			}

			data <- fileInfo
		}
	}()

	return data
}

func GetSumMyself(chanIn <-chan FileInfoMysqlf) <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	go func() {
		defer close(data)

		for fileInfo := range chanIn {
			md5Sum := md5.Sum(fileInfo.Content)          // cari md5 sum
			fileInfo.Sum = hex.EncodeToString(md5Sum[:]) // masukkan md5 sum ke field sum
			data <- fileInfo
		}
	}()

	return data
}

func RenameMyself(chanIn <-chan FileInfoMysqlf) <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	go func() {
		defer close(data)

		for fileInfo := range chanIn {
			newPath := filepath.Join(tempPath, fmt.Sprintf("file-%s.txt.gz", fileInfo.Sum))

			err := os.Rename(fileInfo.FilePath, newPath)

			fileInfo.IsRename = (err == nil)
		}
	}()

	return data
}

func MergeManyChannelMyself(manyChann ...<-chan FileInfoMysqlf) <-chan FileInfoMysqlf {
	data := make(chan FileInfoMysqlf)

	var wg sync.WaitGroup

	wg.Add(len(manyChann))
	for _, eachChan := range manyChann {

		go func(eachChan <-chan FileInfoMysqlf) {
			defer wg.Done()

			for fileInfo := range eachChan {
				data <- fileInfo
			}
		}(eachChan)

	}

	go func() {
		wg.Wait()
		close(data)
	}()

	return data
}
