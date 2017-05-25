package main

import (
	"log"
	"net/http"

	"github.com/alicebob/ssp/ssp"
)

const listen = "localhost:9990"

func main() {
	log.Printf("BidURL: http://%s/rtb", listen)
	log.Fatal(http.ListenAndServe(listen, ssp.Mux()))
}
