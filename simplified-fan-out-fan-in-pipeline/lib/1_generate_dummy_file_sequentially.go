package lib

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

func GenerateFileSequentinally() {
	os.RemoveAll(tempPath)
	os.MkdirAll(tempPath, os.ModePerm)

	for i := 0; i < totalFile; i++ {
		fileName := filepath.Join(tempPath, fmt.Sprintf("file-%d.txt", i))
		content := randomString(contentLength)

		err := os.WriteFile(fileName, []byte(content), os.ModePerm)
		if err != nil {
			log.Printf("Error writting filename %s\n", fileName)
			log.Println(err.Error())
		}

		// setiap satu kali sukses membuat file, logging dijalankan
		// log.Println(i, "file created")
	}

	log.Printf("%d of total file created\n", totalFile)
}

func randomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)

	for index := range b {
		b[index] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
