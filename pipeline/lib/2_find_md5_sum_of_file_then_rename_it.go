// program buat membaca semua file di directory tempPath
// kemudian dicari hash nya
// lalu menggunkan value dari hash tersebut jadi nama baru dari file yang sedang dibaca

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
)

var (
	counterTotal   = 0 // sebagai counter jumplah file yang ditemukan dalam tempPath
	counterRenamed = 0 // sebagai counter jumplah file yang berhasil di rename
)

func Proceed() {

	// filepath.Walk() akan mengiterasi seluruh file dan folder yang sudah dituliskan di parameter pertama
	// setiap file / folder yang ditemukan, maka function walkFunc akan dijalankan
	// parameter kedua dari function filepath.Walk() mempunyai 3 parameter mandatory (path string, info fs.FileInfo, err error)
	err := filepath.Walk(tempPath, walkFunc)

	if err != nil {
		log.Println("ERROR :", err.Error())
	}
	log.Printf("%d/%d file renamed", counterRenamed, counterTotal)
}

// function untuk mengisi parameter kedua function filepath.Walk
func walkFunc(path string, info os.FileInfo, err error) error {
	// jika di setiap pembacaan ditemukan error / directory, maka akan diignore dan dilanjutkan ke iterasi selanjutnya

	// if there is an error, return immediately
	if err != nil {
		return err
	}

	// if it is a sub directory, return immediately
	if info.IsDir() {
		return nil
	}

	counterTotal++

	// read file
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	// compress file
	var dataCompress bytes.Buffer
	gzipCompress, err := gzip.NewWriterLevel(&dataCompress, gzip.BestCompression)
	if err != nil {
		log.Println(err.Error())
	}

	_, err = gzipCompress.Write(buf)
	if err != nil {
		return err
	}

	gzipCompress.Close()

	// write data compressed to a file
	err = os.WriteFile(path, dataCompress.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	// read file lagi
	bfr, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// sum it
	md5 := md5.Sum(bfr)
	sum := hex.EncodeToString(md5[:])

	// rename file
	destinationPath := filepath.Join(tempPath, fmt.Sprintf("file-%s.txt.gz", sum))
	err = os.Rename(path, destinationPath)
	if err != nil {
		return err
	}

	counterRenamed++
	return nil
}
