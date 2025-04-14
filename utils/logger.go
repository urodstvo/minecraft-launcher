package utils

import "log"

type Logger struct{}

func (l *Logger) Info(s string) {
	log.Println(s)
}

func (l *Logger) Error(s string) {
	log.Panic(s)
}