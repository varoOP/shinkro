package logger

import (
	"io"
	"log"
	"os"
)

func NewLogger(path string) *os.File {
	l, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, l)
	log.SetOutput(mw)

	return l
}
