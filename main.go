package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const defaultPort = "8080"

var (
	startTime    time.Time
	requestCount int
)

type Response struct {
	Time         time.Time `json:"time"`
	IP           string    `json:"ip"`
	StartTime    time.Time `json:"startTime"`
	RunningTime  string    `json:"runningTime"`
	RequestCount int       `json:"requestCount"`
	Tag          string    `json:"tag,omitempty"`
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	requestCount += 1
	data := Response{
		Time:         time.Now(),
		IP:           req.RemoteAddr,
		StartTime:    startTime,
		RunningTime:  time.Since(startTime).String(),
		RequestCount: requestCount,
		Tag:          tag(),
	}

	body, _ := json.Marshal(data)
	w.Write(body)
}

func ok(w http.ResponseWriter, req *http.Request) {
	requestCount += 1
	w.Write("OK")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func port() string {
	p := os.Getenv("PORT")
	if len(p) == 0 {
		p = defaultPort
	}
	portNumber, err := strconv.Atoi(p)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Printf("Invalid port number '%d' - Using %s - %v", portNumber, defaultPort, err)
		p = defaultPort
	}
	return p
}

func tag() string {
	return os.Getenv("TAG")
}

func main() {
	startTime = time.Now()
	http.HandleFunc("/", hello)
	http.HandleFunc("/ok", ok)
	serve := fmt.Sprintf(":%s", port())
	log.Printf("http helloworld starting - listening: %s", serve)
	err := http.ListenAndServe(serve, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
