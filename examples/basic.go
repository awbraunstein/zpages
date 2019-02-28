package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/awbraunstein/zpages"
)

var (
	httpAddr = flag.String("http", ":8080", "Listen address")
)

func main() {
	flag.Parse()
	mux := http.NewServeMux()
	requestzHandler, err := zpages.NewRequestz()
	if err != nil {
		log.Fatalf("Unable to initialize Requestz handler; err=%v", err)
	}
	mux.Handle("/healthz", requestzHandler.Middleware(zpages.NewHealthz().ServeHTTP))
	statuszHandler, err := zpages.NewStatusz()
	if err != nil {
		log.Fatalf("Unable to initialize Statusz handler; err=%v", err)
	}
	mux.Handle("/statusz", requestzHandler.Middleware(statuszHandler.ServeHTTP))
	mux.Handle("/requestz", requestzHandler.Middleware(requestzHandler.ServeHTTP))
	log.Println("Listening on %s...", *httpAddr)
	http.ListenAndServe(*httpAddr, mux)
}
