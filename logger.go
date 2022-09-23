package main

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func LoggerInit() {
	file, err := os.OpenFile(*LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal("Failed to read log file")
	}
	InfoLogger = log.New(io.MultiWriter(os.Stdout, file), "[INFO] ", log.LstdFlags|log.Lshortfile)
	WarnLogger = log.New(io.MultiWriter(os.Stdout, file), "[WARNING] ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(io.MultiWriter(os.Stdout, file), "[ERROR] ", log.LstdFlags|log.Lshortfile)
}
