package lib

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	totalFile     = 3000
	contentLength = 500000
)

func GenerateFile() {
	// berbeda dengan os.Remove, os.Remove digunakan untuk menghapus 1 file tunggal / 1 directory yang empty
	// jika digunakan untuk menghapus directory yang ada isinya, maka akan memunculkan error bahwa directory tidak kosong
	// os.RemoveAll digunain buat menghapus path dan semua yang ada didalamnya
	os.RemoveAll(tempPath)

	// berbeda dengan os.Mkdir, jika os.Mkdir digunakan untuk membuat multidirectory dimana directory induk belum ada, maka akan keluar error
	// os.MkdirAll akan membuat semua directory termasuk jika directory induk belum ada, maka akan dibuat dulu
	os.MkdirAll(tempPath, os.ModePerm)

	for i := 0; i < totalFile; i++ {
		fileName := filepath.Join(tempPath, fmt.Sprintf("file-%d.txt", i))
		content := randomString(contentLength) + fmt.Sprint(i) // saya menambahkan statement fmt.Sprint untuk menjaga md5 sum dari tiap file tetap berbeda

		// WriteFile akan menulis data []byte(content) ke fileName,
		// jika fileName belum ada, maka akan dibuat terlebih dahulu dengan permission os.ModePerm
		err := os.WriteFile(fileName, []byte(content), os.ModePerm)
		if err != nil {
			log.Println("Error writting file", fileName, err.Error())
		}

		// logging untuk setiap 100 file dibuat / di generate
		if i%100 == 0 && i > 0 {
			log.Println(i, "file created")
		}
	}

	log.Printf("%d of total file created", totalFile)
}

// randomString digunakan untuk generate data random text sepanjang berdasarkan panjang length
func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for index := range b {
		b[index] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
