package lib

import "testing"

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = randomString(contentLength)
	}
}
