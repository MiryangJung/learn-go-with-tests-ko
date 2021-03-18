package racer

import (
	"net/http"
)

// Racer는 a와 b의 응답 시간을 비교하여 가장 빠른 응답 시간을 반환한다.
func Racer(a, b string) (winner string) {
	select {
	case <-ping(a):
		return a
	case <-ping(b):
		return b
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}
