package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// for exposing pprof functionality
	// to the application
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	NumGoroutines = 1_000
)

func main() {
	// Create lots of goroutines to observe the go runtime scheduling
	// Each goroutine should run in a tight loop so that the scheduler
	// would either preempt or the goroutine yields itself.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Printf("exiting goroutinecreator")
		os.Exit(0)
	}()
	go func() {
		log.Printf("starting profile server")
		_ = http.ListenAndServe("localhost:8080", nil)
	}()
	go func() {
		log.Printf("starting metrics server")
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8090", nil)
	}()
	for i := 0; i < NumGoroutines; i++ {
		if i > 0 && i%10_000 == 0 {
			log.Printf("created %d goroutines", i)
		}
		go spinForever()
	}
	x := 0
	for {
		x++
	}
}

// go:noinline
func spinForever() {
	i := 0
	for {
		i++
		time.Sleep(1 * time.Second)
	}
}
