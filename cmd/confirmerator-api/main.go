package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bahelit/confirmerator/database/mongodb"
)

var (
	client     *mongo.Client
	listenAddr string
	listenPort string
)

const (
	networkAddress = "LISTEN_ADDRESS"
	networkPort    = "LISTEN_PORT"
)

func init() {
	var statusOK bool

	listenAddr, statusOK = os.LookupEnv(networkAddress)
	if !statusOK {
		listenAddr = ":"
	}

	listenPort, statusOK = os.LookupEnv(networkPort)
	if !statusOK {
		listenPort = "80"
	}
}

func main() {
	var err error
	client, err = mongodb.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Printf("ERROR: failed to disconnect from mongo: %v", err)
		}
	}()

	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	// TODO make more robust, don't just give up.
	log.Fatal(http.ListenAndServe(listenAddr+listenPort, router))
}
