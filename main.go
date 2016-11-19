package main

import (
	"encoding/json"
	"fmt"
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
	subscriber := subscriberFactory("socket").(*SocketWorkSubscriber)

	workEngine := NewWorkEngine(MSG_BUFFER_SIZE)
	workEngine.AttachSubscriber(subscriber)

	router := mux.NewRouter()
	configureApi(router, workEngine)
	runServer(router, subscriber, SERVER_PORT)
}

////////////////////////////////////////////////////////////
// Example stdout subscriber to illustrate how decoupled the
// subscribers actually are from the work engine, and how
// they can perform arbirary processing with the work messages
// as long as the IWorkSubscriber interface is implemented
////////////////////////////////////////////////////////////
func subscriberFactory(sub string) IWorkSubscriber {
	// Totally unncessary, but cool nonetheless
	if sub == "socket" {
		return NewSocketWorkSubscriber()
	} else {
		return &StdoutWorkSubscriber{}
	}
}

type StdoutWorkSubscriber struct {
	msgBuffer <-chan string
}

func (s *StdoutWorkSubscriber) Subscribe(msgBuffer <-chan string) {
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
func RunHandler(w http.ResponseWriter, req *http.Request, engine *WorkEngine) {
	res := make(map[string]string)
	res["job_token"] = engine.StartNewJob(NewFooMachine())
	json.NewEncoder(w).Encode(res)
}

////////////////////////////////////////////////////////////
// Server configuration
////////////////////////////////////////////////////////////
func configureApi(router *mux.Router, workEngine *WorkEngine) {
	// NOTE: These paths must absolutely match the same paths that
	// the router is mounted at!
	// Ex: http.Handle("/api/", r), then anything on r under this space
	// must be `/api/$PATH`
	router.HandleFunc(
		"/api/run", createHandler(workEngine, RunHandler),
	).Methods("POST")
}

// Sets up web server-y things like static and template handlers
func runServer(
	router *mux.Router,
	sws *SocketWorkSubscriber,
	port int,
) {
	http.Handle("/socket.io/", sws.server) // Mount the socket.io server
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.Handle("/api/", router)
	http.HandleFunc("/", IndexHandler) // Mount the lone template server

	portStr := fmt.Sprintf(":%d", port)
	fmt.Println(fmt.Sprintf("Listening on localhost:%d", port))
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, nil)
}

////////////////////////////////////////////////////////////
// Util
////////////////////////////////////////////////////////////
// Want ability to create route handlers that are conformant with vanilla
// gorilla handlers, but have an injected work engine reference via closure
// Desire is to favor dependency injection over package level globals!
type GorillaRouteHandler func(http.ResponseWriter, *http.Request)
type InjectedRouteHandler func(http.ResponseWriter, *http.Request, *WorkEngine)

func createHandler(engine *WorkEngine, r InjectedRouteHandler) GorillaRouteHandler {
	return func(writer http.ResponseWriter, request *http.Request) {
		r(writer, request, engine)
	}
}
