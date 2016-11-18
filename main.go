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
	"time"
)

const MSG_BUFFER_SIZE = 10
const SERVER_PORT = 3000

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // TURBOBOOST!
	rand.Seed(time.Now().UnixNano())
}

func main() {
	router := mux.NewRouter()
	workManager := NewWorkManager(MSG_BUFFER_SIZE)
	workManager.AttachSubscriber(subscriberFactory("stdout"))

	configureApi(router, workManager)
	runServer(router, SERVER_PORT)
}

////////////////////////////////////////////////////////////
// Example stdout subscriber to illustrate how decoupled the
// subscribers actually are from the worker manager, and how
// they can perform arbirary processing with the work messages
// as long as the IWorkSubscriber interface is implemented
////////////////////////////////////////////////////////////
func subscriberFactory(sub string) IWorkSubscriber {
	// Totally unncessary, but cool nonetheless
	if sub == "socket" {
		panic("SOCKETSUBSCRIBER NOT YET IMPLEMENTED")
	} else {
		return &StdoutSubscriber{}
	}
}

type StdoutSubscriber struct {
	msgBuffer <-chan string
}

func (s *StdoutSubscriber) Subscribe(msgBuffer <-chan string) {
	// Always drain the buffer if there's a message waiting.
	// Here we're just forwarding to stdout, but of course, the message
	// destination could be anything (ultimate websockets!)
	// NOTE: DON'T FORGET TO GOROUTINE THIS, OR WILL YOU CHOKE THE MAIN PROCESSOR
	s.msgBuffer = msgBuffer
	go func() {
		for {
			msg := <-msgBuffer
			fmt.Printf(msg)
		}
	}()
}

////////////////////////////////////////////////////////////
// API Handlers
////////////////////////////////////////////////////////////
func RunHandler(w http.ResponseWriter, req *http.Request, wm *WorkManager) {
	res := make(map[string]string)
	res["job_token"] = wm.StartNewJob(NewFooMachine())
	json.NewEncoder(w).Encode(res)
}

////////////////////////////////////////////////////////////
// Server configuration
////////////////////////////////////////////////////////////
func configureApi(router *mux.Router, workManager *WorkManager) {
	router.HandleFunc(
		"/run", createHandler(workManager, RunHandler),
	).Methods("POST")
}

// Sets up web server-y things like static and template handlers
func runServer(router *mux.Router, port int) {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	router.PathPrefix("/static/").Handler(fs)
	router.HandleFunc("/", IndexHandler)

	fmt.Println(fmt.Sprintf("Listening on localhost:%d", port))
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	portStr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(portStr, handlers.CORS(
		allowedHeaders,
	)(router)))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, nil)
}

////////////////////////////////////////////////////////////
// Util
////////////////////////////////////////////////////////////
// Want ability to create route handlers that are conformant with vanilla
// gorilla handlers, but have an injected work manager reference via closure
// Desire is to favor dependency injection over package level globals!
type GorillaRouteHandler func(http.ResponseWriter, *http.Request)
type InjectedRouteHandler func(http.ResponseWriter, *http.Request, *WorkManager)

func createHandler(wm *WorkManager, r InjectedRouteHandler) GorillaRouteHandler {
	return func(writer http.ResponseWriter, request *http.Request) {
		r(writer, request, wm)
	}
}
