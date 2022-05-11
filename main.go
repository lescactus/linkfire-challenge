package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	cfg *Config
)

// App represents the 'payment' application.
// It comes with a custom HTTP server and can be extended
type App struct {
	// HTTP server with custom handler and read and write timeouts
	Server *http.Server
}

// HelloRequest represents the json request struct send to the Hello handler.
type HelloRequest struct {
	Message string `json:"message"`
}

// PingResponse represents the json response to the Ping handler.
// It provides some runtime informations
type PingResponse struct {
	Hostname string `json:"hostname"`
	Message  string `json:"message"`
	GOOS     string `json:"goos"`
	GOARCH   string `json:"goarch"`
	Runtime  string `json:"runtime"`
	NumCPU   int    `json:"cpu"`
	AppName  string `json:"application_name"`
}

// Run will call ListenAndServe and wrap it in a log.Fatal() to exit
// to any error encountered
func (a *App) Run() {
	log.Println("Starting http server on port " + a.Server.Addr)
	// Bind http server configured address and start to serve requests
	log.Fatal(a.Server.ListenAndServe())
}

// healthCheck provide a simple health endpoint, typically useful
// to be used by a Kubernetes readinessProbe/livenessProbe.
//
// ref: https://inadarei.github.io/rfc-healthcheck/#
//
// The above docment provides a good starting point on what to include in a API healthcheck.
// In our case, 'payment' is a very minimal service without external dependencies
// The health check is then a very simple ping-pong.
func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	// Set a 200 OK response with a Content-Type http header
	w.Header().Set("Content-Type", "application/health+json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"status":"pass"}`)
}

// Ping provides a simple GET handler used to send some runtime informations
// to the client.
func (a *App) Ping(w http.ResponseWriter, r *http.Request) {
	p := NewPingResponse("Pong")

	data, err := json.Marshal(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(data))
}

// Hello provides a simple POST handler handler used to send some runtime informations
// to the client. It expects a HelloRequest and will use its Message field to answer with
// a PingResponse
func (a *App) Hello(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var d HelloRequest
	err = json.Unmarshal(body, &d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := NewPingResponse(d.Message)
	data, err := json.Marshal(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(data))
}

// NewApp returns a new app instance.
// It takes in parameters the tcp address for the server to listen on
// and a pointer to a mux router.
func NewApp(addr string, r *mux.Router) *App {
	return &App{
		Server: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadTimeout:       cfg.GetDuration("SERVER_READ_TIMEOUT"),
			ReadHeaderTimeout: cfg.GetDuration("SERVER_READ_HEADER_TIMEOUT"),
			WriteTimeout:      cfg.GetDuration("SERVER_WRITE_TIMEOUT"),
		},
	}
}

// NewPingResponse returns a new PingResponse instance.
// It will read the runtime informations from the 'runtime' package.
// It takes in argument a message string
func NewPingResponse(msg string) *PingResponse {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "<undefined>"
	}

	return &PingResponse{
		Hostname: hostname,
		Message:  msg,
		GOOS:     runtime.GOOS,
		GOARCH:   runtime.GOARCH,
		Runtime:  runtime.Version(),
		NumCPU:   runtime.NumCPU(),
		AppName:  AppName,
	}
}

func init() {
	cfg = NewConfig()
}

func main() {
	r := mux.NewRouter()
	app := NewApp(cfg.GetString("APP_ADDR"), r) // using port tcp/8080 by default

	// Routes registration
	r.HandleFunc("/rest/ready", app.healthCheck).Methods("GET")
	r.HandleFunc("/rest/alive", app.healthCheck).Methods("GET")
	r.HandleFunc("/rest/v1/ping", app.Ping).Methods("GET")
	r.HandleFunc("/rest/v1/hello", app.Hello).Methods("POST")

	// Using logging middleware
	r.Use(func(next http.Handler) http.Handler { return handlers.LoggingHandler(os.Stdout, next) })

	app.Run()
}
