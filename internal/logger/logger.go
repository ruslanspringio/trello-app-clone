package logger

import (
	"log"
	"time"
)

var (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

func LogRequest(status int, method, path, clientIP string, latency time.Duration) {
	color := green
	switch {
	case status >= 500:
		color = red
	case status >= 400:
		color = yellow
	}

	log.Printf("%s[%d]%s %s %s %s (%v)",
		color, status, reset,
		method,
		path,
		clientIP,
		latency,
	)
}
