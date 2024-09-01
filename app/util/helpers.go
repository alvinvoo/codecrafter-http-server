package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
)

// Debug logger function
func DebugLog(title string, message ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		log.Println("DEBUG:", title, message)
	}
}

func GzipEncode(body string) (string, error) {
	// Implement gzip encoding here
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(body))
	if err != nil {
		return "", fmt.Errorf("error writing compressed data: %s", err.Error())
	}

	// must close to flush remaining compressed data to the buffer
	if err := zw.Close(); err != nil {
		return "", fmt.Errorf("error writing compressed data: %s", err.Error())
	}

	return buf.String(), nil
}
