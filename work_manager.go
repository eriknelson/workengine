package main

import (
	"github.com/ventu-io/go-shortid"
)

type IWork interface {
	Run(token string, msgBuffer chan<- string)
}

type IWorkSubscriber interface {
	Subscribe(msgBuffer <-chan string)
}

type WorkManager struct {
	msgBuffer chan string
}

func NewWorkManager(bufferSize int) *WorkManager {
	return &WorkManager{
		msgBuffer: make(chan string, bufferSize),
	}
}

func (m *WorkManager) StartNewJob(work IWork) string {
	jobToken, _ := shortid.Generate()
	go work.Run(jobToken, m.msgBuffer)
	return jobToken
}

func (m *WorkManager) AttachSubscriber(subscriber IWorkSubscriber) {
	subscriber.Subscribe(m.msgBuffer)
}
