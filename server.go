package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"
)

type key int

const (
	requestIDKey key = 0
)

var (
	logger     = log.New(os.Stdout, "", log.LstdFlags)
	healthy    int32
	listenAddr = ":" + port()
)

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

func createServer() *http.Server {

	logger.Println("Server is starting...")

	router := http.NewServeMux()
	routes := SetupHttpCustomHandlersRouter(router)
	handledRoot := false
	handledOk := false
	for _, route := range routes {
		switch route {
		case "/":
			handledRoot = true
		case "/ok":
			handledOk = true
		}
	}
	if !handledRoot {
		router.HandleFunc("/", hello)
	}
	if !handledOk {
		router.HandleFunc("/ok", ok)
	}
	for _, route := range GetRoutesRouter(router) {
		log.Printf("Route %s", route)
	}
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port()),
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return server
}

func runServer(server *http.Server) {
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			o := &responseObserver{ResponseWriter: w}

			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, o.status, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(o, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			rqHeader := fmt.Sprintf("%d", requestCount.Value())
			w.Header().Set("X-Request-Count", rqHeader)	// TODO: request count not working on response
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
