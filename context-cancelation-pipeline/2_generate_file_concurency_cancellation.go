package main

import (
	"contex-cancelation-pipeline/lib"
	"context"
	"log"
	"time"
)

func main() {
	log.Println("start")
	start := time.Now()

	// WithTimeout digunakan untuk menambahkan timeout pada sebuah contex
	ctx, cancel := context.WithTimeout(context.Background(), lib.TimeDuration)
	// ctx, cancel := context.WithCancel(context.Background())

	// meskipun contex sudah punya timeout dan kita tidak perlu meng-cancle secara manual,
	// sangat dianjurkan untuk tetap memanggil callback cancle() secara deffered
	defer cancel()

	// cara pembuatan object contex
	//
	//
	// 1. menggunakan function context.Background, menghasilkan object contex yang data didalamnya adalah kosong dan tidak memiliki deadline.
	// biasanya digunakan untuk inisialisasi object context baru yang nantinya akan dichain dengan function context.With...
	//
	// 2. menggunakan function context.TODO, menghasilkan object context baru sama seperti context.Background,
	// namun ini biasanya digunakan dalam situasi ketika belum jelas nantinya harus menggunakan contect apa (apakah dengan timeout apa cancle)
	//
	// 3. menggunakan function context.With.., sebenarnya bukan untuk menghasilkan object context baru,
	// tapi digunakan untuk menambahkan informasi pada copied context yang telah disisipkan pada parameter pertama
	//
	//
	// ada 3 buah function context.With.. yang bisa digunakan :
	//
	// 1. context.WithCancle(ctx) (ctx, cancel) menambahkan fasilitas cancellable pada context yang disisipkan pada parameter pertama
	// nilai balik dari function kedua ini adalah cancel yang tipenya context.CancelFunc kita bisa secara paksa mencancel context ini
	//
	// 2. contex.WithDeadline(ctx, time.Time) (ctx, cancel) juga menambahkan fasilitas cancellable pada context yang disisipkan,
	// namun juga menambahkan informasi deadline yang dimana jika waktu sekarang sudah melebihi deadline yang sudah ditentukan, maka context secara otomatis dicancel secara paksa.
	//
	// 3. context.WithTimeout(ctx, time.Duration) (ctx, cancel) sama seperti contex.WithDeadline bedanya parameter kedua bertipe time.Duration

	// GenerateFilesWithContext function yang sama persis dengan GenerateFiles
	// perbedaanya pada function ini diterapkan calcellation dengan Context
	lib.GenerateFilesWithContext(ctx)

	duration := time.Since(start)
	log.Printf("done in %f seconds\n", duration.Seconds())
}
