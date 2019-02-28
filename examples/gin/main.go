package main

import (
	"flag"
	"log"

	"github.com/awbraunstein/zpages"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
)

var (
	httpAddr = flag.String("http", ":8080", "Listen address")
)

func main() {
	flag.Parse()
	r := gin.New()
	requestzHandler, err := zpages.NewRequestz()
	if err != nil {
		log.Fatalf("Unable to initialize Requestz handler; err=%v", err)
	}
	r.Use(adapter.Wrap(requestzHandler.Middleware))

	r.GET("/healthz", gin.WrapH(zpages.NewHealthz()))
	statuszHandler, err := zpages.NewStatusz()
	if err != nil {
		log.Fatalf("Unable to initialize Statusz handler; err=%v", err)
	}
	r.GET("/statusz", gin.WrapH(statuszHandler))
	r.GET("/requestz", gin.WrapH(requestzHandler))
	log.Println("Listening on %s...", *httpAddr)
	r.Run(*httpAddr)
}
