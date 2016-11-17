package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	INITIALIZING = iota
	RUNNING      = iota
	FINALIZING   = iota
	FINISHED     = iota
)

var (
	stateText map[int]string
	wg        sync.WaitGroup
)

type WorkMachine struct {
	identifier    string
	stateGraph    map[int]int
	currentState  int
	stateWeights  map[int]int
	totalWeight   int
	totalTicks    int
	stateTicks    int
	totalProgress float32
	stateProgress float32
}

func NewWorkMachine(identifier string) WorkMachine {
	stateGraph := map[int]int{
		INITIALIZING: RUNNING,
		RUNNING:      FINALIZING,
		FINALIZING:   FINISHED,
	}

	stateWeights := make(map[int]int)

	shortWeight := random(3, 7)
	longWeight := random(7, 12)
	totalWeight := shortWeight*2 + longWeight

	stateWeights[INITIALIZING] = shortWeight
	stateWeights[RUNNING] = longWeight
	stateWeights[FINALIZING] = shortWeight

	return WorkMachine{
		identifier:   identifier,
		stateWeights: stateWeights,
		stateGraph:   stateGraph,
		totalWeight:  totalWeight,
	}
}

func (m *WorkMachine) Tick() {
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

func (m *WorkMachine) Report() {
	fmt.Printf("%s -> State: %s, State Progress: %.2f, Total Progress: %.2f \n",
		m.identifier, stateText[m.currentState], m.stateProgress, m.totalProgress)
}

func runMachine(identifier string) {
	machine := NewWorkMachine(identifier)

	go func() {
		fmt.Printf("total weight -> %d\n", machine.totalWeight)
		fmt.Println(machine.stateWeights)

		for {
			if machine.currentState == FINISHED {
				break
			}

			machine.Tick()
			if machine.currentState != FINISHED {
				machine.Report()
			}
			time.Sleep(time.Millisecond * 500)
		}

		fmt.Printf("Finished -> %d, %d, %.2f\n",
			machine.currentState,
			machine.totalTicks,
			machine.totalProgress,
		)
	}()
}

func init() {
	runtime.GOMAXPROCS(2)

	stateText = map[int]string{
		INITIALIZING: "Initializing",
		RUNNING:      "Running",
		FINALIZING:   "Finalizing",
		FINISHED:     "Finished",
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/run", RunHandler).Methods("POST")

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	r.PathPrefix("/static/").Handler(fs)
	r.HandleFunc("/", IndexHandler)

	fmt.Println("Listening on localhost:3000")
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(
		allowedHeaders,
	)(r)))
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&body)
	identifier := body["id"]

	runMachine(identifier)

	resMap := map[string]string{"foo": "bar"}
	json.NewEncoder(w).Encode(resMap)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, nil)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
