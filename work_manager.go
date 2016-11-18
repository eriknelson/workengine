package main

import (
	"github.com/ventu-io/go-shortid"
)

type IWork interface {
	Run(token string, msgBuffer chan<- string)
}

type WorkManager struct {
	msgBuffer chan string
}

func NewWorkManager(bufferSize int) *WorkManager {
	return &WorkManager{
		msgBuffer: make(chan string, bufferSize),
	}
}

// Returns job token
func (m *WorkManager) StartNewJob(work IWork) string {
	jobToken, _ := shortid.Generate()
	go work.Run(jobToken, m.msgBuffer)
	return jobToken
}
