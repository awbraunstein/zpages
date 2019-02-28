package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/awbraunstein/zpages"
	"github.com/gorilla/mux"
)

var (
	httpAddr = flag.String("http", ":8080", "Listen address")
)

func main() {
	flag.Parse()
	r := mux.NewRouter()
	requestzHandler, err := zpages.NewRequestz()
	if err != nil {
		log.Fatalf("Unable to initialize Requestz handler; err=%v", err)
	}
	r.Use(requestzHandler.Middleware)
	r.Handle("/healthz", zpages.NewHealthz())
	statuszHandler, err := zpages.NewStatusz()
	if err != nil {
		log.Fatalf("Unable to initialize Statusz handler; err=%v", err)
	}
	r.Handle("/statusz", statuszHandler)
	r.Handle("/requestz", requestzHandler)
	log.Printf("Listening on %s...", *httpAddr)
	http.ListenAndServe(*httpAddr, r)
}
