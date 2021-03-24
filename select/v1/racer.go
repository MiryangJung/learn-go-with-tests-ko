package racer

import (
	"net/http"
	"time"
)

// Racer는 a와 b의 응답 시간을 비교하여 가장 빠른 응답 시간을 반환한다.
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
