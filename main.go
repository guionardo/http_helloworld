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

var startTime time.Time

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	data := map[string]string{
		"time":        time.Now().String(),
		"ip":          req.RemoteAddr,
		"startTime":   startTime.String(),
		"runningTime": time.Since(startTime).String(),
	}

	body, _ := json.Marshal(data)
	w.Write(body)
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
	if err != nil {
		log.Printf("Invalid port number '%s' - Using %s", p, defaultPort)
	} else {
		if portNumber < 1 || portNumber > 65535 {
			log.Printf("Invalid port number '%d' - Using %s", portNumber, defaultPort)
			p = defaultPort
		} else {
			p = fmt.Sprintf("%d", portNumber)
		}
	}
	return p
}

func main() {
	startTime = time.Now()
	http.HandleFunc("/", hello)
	serve := fmt.Sprintf(":%s", port())
	log.Printf("http helloworld starting - listening: %s", serve)
	err := http.ListenAndServe(serve, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
