package main

import (
	//"encoding/json"
	"fmt"
	//"github.com/gorilla/handlers"
	//"github.com/gorilla/mux"
	//"html/template"
	//"log"
	"math/rand"
	//"net/http"
	"runtime"
	"strings"
	"time"
)

////////////////////////////////////////////////////////////
// PARAM CONFIG
const WORKER_COUNT = 4
const MAX_PROC_COUNT = 4
const CLOCK_PERIOD = 500

////////////////////////////////////////////////////////////

func init() {
	runtime.GOMAXPROCS(MAX_PROC_COUNT)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	msgBuffer := make(chan string)

	for i := 0; i < WORKER_COUNT; i++ {
		NewFooMachine(msgBuffer).Run()
	}

	finishedCount := 0
	// TODO: Simulating a MessageConsumer pulling worker messages off the queue
	// and processing them (we're simply writing to stdout)
	for {
		msg := <-msgBuffer
		fmt.Printf(msg)

		if strings.Contains(msg, "Finished") {
			finishedCount++
			if finishedCount == WORKER_COUNT {
				break
			}
		}
	}

	fmt.Println("FIN")
}

//r := mux.NewRouter()
//r.HandleFunc("/run", RunHandler).Methods("POST")

//fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
//r.PathPrefix("/static/").Handler(fs)
//r.HandleFunc("/", IndexHandler)

//fmt.Println("Listening on localhost:3000")
//allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
//log.Fatal(http.ListenAndServe(":3000", handlers.CORS(
//allowedHeaders,
//)(r)))

//func RunHandler(w http.ResponseWriter, r *http.Request) {
//body := make(map[string]string)
//json.NewDecoder(r.Body).Decode(&body)
//identifier := body["id"]

//runMachine(identifier)

//resMap := map[string]string{"foo": "bar"}
//json.NewEncoder(w).Encode(resMap)
//}

//func IndexHandler(w http.ResponseWriter, r *http.Request) {

//t, _ := template.ParseFiles("static/index.html")
//t.Execute(w, nil)
//}
