package main

import (
	"log"
	"net/http"

	"github.com/alicebob/ssp/ssp"
)

const (
	listen = ":9998"
)

var (
	dsps = []ssp.DSP{
		{
			ID:     "dsp1",
			Name:   "Test 1",
			BidURL: "http://localhost:9990/rtb",
		},
		{
			ID:     "dsp2",
			Name:   "Test 2 - offline",
			BidURL: "http://localhost:9999/...",
		},
	}
	placements = []ssp.Placement{
		ssp.Placement{
			ID:     "my_website_1",
			Name:   "My Website",
			Width:  520,
			Height: 100,
		},
		ssp.Placement{
			ID:     "my_website_2",
			Name:   "My Website 2",
			Width:  300,
			Height: 330,
		},
	}
)

func main() {
	s := NewDaemon(dsps)
	log.Printf("listening on %s...", listen)
	log.Fatal(http.ListenAndServe(listen, mux(s, placements)))
}
