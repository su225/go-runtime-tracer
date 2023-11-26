package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// for exposing pprof functionality
	// to the application
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	NumGoroutines = 10_000_000
)

func main() {
	// Create 10 million goroutines to observe the go runtime scheduling
	// Each goroutine should run in a tight loop so that the scheduler
	// would either preempt or the goroutine yields itself.
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go log.Println(http.ListenAndServe("localhost:8080", nil))
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8080", nil)
	}()
	for i := 0; i < NumGoroutines; i++ {
		go spinForever()
	}
	<-sigChan
}

func spinForever() {
	i := 0
	for {
		i++
	}
}
