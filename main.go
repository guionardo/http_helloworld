package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = "8080"
	version     = "0.9.4"
	toolName    = "http-helloworld"
)

var (
	startTime    time.Time
	requestCount = &RequestCount{}
)

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	requestCount.Inc()
	data := Response{
		Time:         time.Now(),
		IP:           req.RemoteAddr,
		StartTime:    startTime,
		RunningTime:  time.Since(startTime).String(),
		RequestCount: requestCount.Value(),
		Tag:          os.Getenv("TAG"),
	}

	body, _ := json.Marshal(data)
	w.Write(body)
}

func ok(w http.ResponseWriter, req *http.Request) {
	requestCount.Inc()
	fmt.Fprintf(w, "OK")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	log.Printf("%s v%s", toolName, version)
	startTime = time.Now()

	server := createServer()
	runServer(server)

}
