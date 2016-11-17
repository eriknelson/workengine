package main

import (
	"fmt"
	"math/rand"
	"time"
)

var FOO_ID int

func init() {
	FOO_ID = 0
}

const (
	FOO_INITIALIZING = iota
	FOO_RUNNING      = iota
	FOO_FINALIZING   = iota
	FOO_FINISHED     = iota
)

type FooMachine struct {
	identifier    int
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

func NewFooMachine(msgBuffer chan<- string) *FooMachine {
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

	FOO_ID++

	return &FooMachine{
		identifier:   FOO_ID,
		stateWeights: stateWeights,
		stateGraph:   stateGraph,
		stateText:    stateText,
		totalWeight:  totalWeight,
		msgBuffer:    msgBuffer,
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
	////////////////////////////////////////////////////////////
	// TOOD: Fix edge case for progress
	var stateProgress, totalProgress float32
	if m.currentState == FOO_FINISHED {
		stateProgress = 1.00
		totalProgress = 1.00
	} else {
		stateProgress = m.stateProgress
		totalProgress = m.totalProgress
	}
	////////////////////////////////////////////////////////////

	m.emit(
		fmt.Sprintf("%d -> State: %s, State Progress: %.2f, Total Progress: %.2f \n",
			m.identifier, m.stateText[m.currentState], stateProgress, totalProgress),
	)
}

func (m *FooMachine) Run() {
	go func() {
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
	}()
}

func (m *FooMachine) emit(msg string) {
	m.msgBuffer <- msg
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
