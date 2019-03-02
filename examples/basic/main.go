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
	mux.Handle("/healthz", requestzHandler.Middleware(zpages.NewHealthz()))
	statuszHandler, err := zpages.NewStatusz(func(data map[string]string) {
		data["Hello Message"] = "This is a Statusz page."

	})
	if err != nil {
		log.Fatalf("Unable to initialize Statusz handler; err=%v", err)
	}
	mux.Handle("/statusz", requestzHandler.Middleware(statuszHandler))
	mux.Handle("/requestz", requestzHandler.Middleware(requestzHandler))
	log.Printf("Listening on %s...", *httpAddr)
	http.ListenAndServe(*httpAddr, mux)
}
