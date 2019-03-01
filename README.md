# zpages
[![MIT licensed](https://img.shields.io/github/license/awbraunstein/zpages.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/awbraunstein/zpages?status.svg)](https://godoc.org/github.com/awbraunstein/zpages)
[![Build Status](https://travis-ci.org/awbraunstein/zpages.svg?branch=master)](https://travis-ci.org/awbraunstein/zpages)
[![Coverage Status](https://img.shields.io/codecov/c/github/awbraunstein/zpages.svg)](https://codecov.io/gh/awbraunstein/zpages)

Go utilities for generating helpful debug and internal pages for server inspection.

## Installation

`go get -u github.com/awbraunstein/zpages`

## Usage

Just add the handlers/middleware that you want and run your server. It is recommended to put these pages behind some sort of internal auth to avoid leaking this information to users.
```golang
func main() {
	mux := http.NewServeMux()
	requestzHandler, err := zpages.NewRequestz()
	if err != nil {
		log.Fatalf("Unable to initialize Requestz handler; err=%v", err)
	}
	mux.Handle("/healthz", requestzHandler.Middleware(zpages.NewHealthz()))
	statuszHandler, err := zpages.NewStatusz()
	if err != nil {
		log.Fatalf("Unable to initialize Statusz handler; err=%v", err)
	}
	mux.Handle("/statusz", requestzHandler.Middleware(statuszHandler))
	mux.Handle("/requestz", requestzHandler.Middleware(requestzHandler))
	log.Println("Listening on %s...", *httpAddr)
	http.ListenAndServe(":8080", mux)
}
```
