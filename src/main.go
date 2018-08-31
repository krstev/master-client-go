package main

import (
	"log"
	"master-client-go/src/eureka"
	"master-client-go/src/service"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func main() {
	handleSigterm() // Graceful shutdown on Ctrl+C or kill

	//go startWebServer() // Starts HTTP service  (async)
	portC := make(chan int, 1)
	go startWebServer(portC) // Starts HTTP service  (async)

	port := <-portC
	close(portC)
	eureka.Register(port) // Performs Eureka registration

	go eureka.StartHeartbeat() // Performs Eureka heartbeating (async)

	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block goclient() exit
	wg.Add(1)
	wg.Wait()
}

func handleSigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		eureka.Deregister()
		os.Exit(1)
	}()
}

func startWebServer(port chan int) int {
	router := service.NewRouter()
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port <- listener.Addr().(*net.TCPAddr).Port
	log.Println("Starting HTTP service at ", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, router))
	//err1 := http.ListenAndServe(":5555", router)
	//if err1 != nil {
	//	log.Println("An error occured starting HTTP listener at port 5555")
	//	log.Println("Error: " + err.Error())
	//}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contetn-Type", "application/json")
	bs := []byte(strconv.Itoa(rand.Intn(100)))
	w.Write(bs)
}
