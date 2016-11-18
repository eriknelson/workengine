package main

import (
	"fmt"
	"math/rand"
	"time"
)

const CLOCK_PERIOD = 500

const (
	FOO_INITIALIZING = iota
	FOO_RUNNING      = iota
	FOO_FINALIZING   = iota
	FOO_FINISHED     = iota
)

type FooMachine struct {
	jobToken      string
	stateGraph    map[int]int
	stateText     map[int]string
	currentState  int
	stateWeights  map[int]int
	totalWeight   int
	totalTicks    int
	stateTicks    int
	totalProgress float32
	stateProgress float32
	msgBuffer     chan<- string
}

func NewFooMachine() *FooMachine {
	stateGraph := map[int]int{
		FOO_INITIALIZING: FOO_RUNNING,
		FOO_RUNNING:      FOO_FINALIZING,
		FOO_FINALIZING:   FOO_FINISHED,
	}

	stateText := map[int]string{
		FOO_INITIALIZING: "Initializing",
		FOO_RUNNING:      "Running",
		FOO_FINALIZING:   "Finalizing",
		FOO_FINISHED:     "Finished",
	}

	stateWeights := make(map[int]int)

	shortWeight := random(3, 7)
	longWeight := random(7, 12)
	totalWeight := shortWeight*2 + longWeight

	stateWeights[FOO_INITIALIZING] = shortWeight
	stateWeights[FOO_RUNNING] = longWeight
	stateWeights[FOO_FINALIZING] = shortWeight

	return &FooMachine{
		stateWeights: stateWeights,
		stateGraph:   stateGraph,
		stateText:    stateText,
		totalWeight:  totalWeight,
	}
}

func (m *FooMachine) tick() {
	// check to see if we need to transition to the next state
	if m.stateWeights[m.currentState] == m.stateTicks {
		m.stateTicks = 0                              // Reset state ticker
		m.currentState = m.stateGraph[m.currentState] // Transition to next state
	}

	m.totalTicks++
	m.stateTicks++

	m.totalProgress = float32(m.totalTicks) / float32(m.totalWeight)
	m.stateProgress = float32(m.stateTicks) / float32(m.stateWeights[m.currentState])
}

func (m *FooMachine) report() {
	//////////////////////////////////////////////////////////
	//TOOD: Fix edge case for progress
	var stateProgress, totalProgress float32
	if m.currentState == FOO_FINISHED {
		stateProgress = 1.00
		totalProgress = 1.00
	} else {
		stateProgress = m.stateProgress
		totalProgress = m.totalProgress
	}
	//////////////////////////////////////////////////////////

	m.emit(
		fmt.Sprintf("State: %s, State Progress: %.2f, Total Progress: %.2f \n",
			m.stateText[m.currentState], stateProgress, totalProgress),
	)
}

func (m *FooMachine) Run(jobToken string, msgBuffer chan<- string) {
	m.msgBuffer = msgBuffer
	m.jobToken = jobToken

	m.emit(fmt.Sprintf("total weight -> %d\n", m.totalWeight))
	m.emit(fmt.Sprintln(m.stateWeights))

	for {
		if m.currentState == FOO_FINISHED {
			break
		}

		m.tick()
		m.report()
		time.Sleep(time.Millisecond * CLOCK_PERIOD)
	}
}

func (m *FooMachine) emit(msg string) {
	m.msgBuffer <- fmt.Sprintf("[%s] %s", m.jobToken, msg)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
