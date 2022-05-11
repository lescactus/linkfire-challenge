package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	r := mux.NewRouter()
	app := NewApp(":8080", r)

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/rest/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder satisfying http.ResponseWriter to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.healthCheck)
	handler.ServeHTTP(rr, req)

	expectedHTTPStatus := http.StatusOK
	expectedContentType := "application/health+json"
	expectedBody := `{"status":"pass"}`

	assert.Equal(t, expectedHTTPStatus, rr.Code)
	assert.Equal(t, expectedContentType, rr.Result().Header.Get("Content-Type"))
	assert.Equal(t, expectedBody, rr.Body.String())
}

func TestPing(t *testing.T) {
	r := mux.NewRouter()
	app := NewApp(":8080", r)

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/rest/v1/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder satisfying http.ResponseWriter to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.Ping)
	handler.ServeHTTP(rr, req)

	p := NewPingResponse("Pong")
	expectedHTTPStatus := http.StatusOK
	expectedContentType := "application/json"
	expectedBody, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedHTTPStatus, rr.Code)
	assert.Equal(t, expectedContentType, rr.Result().Header.Get("Content-Type"))
	assert.Equal(t, string(expectedBody), rr.Body.String())
}

func TestHello(t *testing.T) {
	r := mux.NewRouter()
	app := NewApp(":8080", r)

	tests := []struct {
		name                string
		app                 *App
		payload             interface{}
		expectedHTTPStatus  int
		expectedContentType string
	}{
		{
			name:                "Valid HelloRequest json - empty message",
			app:                 app,
			payload:             &HelloRequest{Message: ""},
			expectedHTTPStatus:  http.StatusOK,
			expectedContentType: "application/json",
		},
		{
			name:                "Valid HelloRequest json - non empty message",
			app:                 app,
			payload:             &HelloRequest{Message: "Hello world"},
			expectedHTTPStatus:  http.StatusOK,
			expectedContentType: "application/json",
		},
		{
			name:                "Invalid json - non empty message",
			app:                 app,
			payload:             `invalid payload`,
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedContentType: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(test.payload)
			if err != nil {
				t.Fatal(err)
			}

			// Create a request to pass to the handler
			req, err := http.NewRequest("POST", "/rest/v1/hello", &buf)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder satisfying http.ResponseWriter to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.Hello)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedHTTPStatus, rr.Code)
			assert.Equal(t, test.expectedContentType, rr.Result().Header.Get("Content-Type"))
		})
	}
}

func TestNewApp(t *testing.T) {
	type args struct {
		addr string
		r    *mux.Router
	}
	tests := []struct {
		name string
		args args
		want *App
	}{
		{
			name: "Addr = ':8080'",
			args: args{addr: ":8080", r: mux.NewRouter()},
			want: &App{Server: &http.Server{Addr: ":8080", Handler: mux.NewRouter(), ReadTimeout: 30 * time.Second, ReadHeaderTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second}},
		},
		{
			name: "Addr = '127.0.0.1:8080'",
			args: args{addr: "127.0.0.1:8080", r: mux.NewRouter()},
			want: &App{Server: &http.Server{Addr: "127.0.0.1:8080", Handler: mux.NewRouter(), ReadTimeout: 30 * time.Second, ReadHeaderTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second}},
		},
		{
			name: "Addr = '0.0.0.0:8080'",
			args: args{addr: "0.0.0.0:8080", r: mux.NewRouter()},
			want: &App{Server: &http.Server{Addr: "0.0.0.0:8080", Handler: mux.NewRouter(), ReadTimeout: 30 * time.Second, ReadHeaderTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.args.addr, tt.args.r)
			assert.Equal(t, tt.want, app)
		})
	}
}

func TestNewPingResponse(t *testing.T) {
	const (
		undefinedHostnameValue = "<undefined>"
	)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = undefinedHostnameValue
	}

	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want *PingResponse
	}{
		{
			name: "Empty message",
			args: args{msg: ""},
			want: &PingResponse{
				Hostname: hostname,
				Message:  "",
				GOOS:     runtime.GOOS,
				GOARCH:   runtime.GOARCH,
				Runtime:  runtime.Version(),
				NumCPU:   runtime.NumCPU(),
				AppName:  AppName,
			},
		},
		{
			name: "Non empty message",
			args: args{msg: "hello world"},
			want: &PingResponse{
				Hostname: hostname,
				Message:  "hello world",
				GOOS:     runtime.GOOS,
				GOARCH:   runtime.GOARCH,
				Runtime:  runtime.Version(),
				NumCPU:   runtime.NumCPU(),
				AppName:  AppName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPingResponse(tt.args.msg)
			assert.Equal(t, tt.want, got)
		})
	}
}
