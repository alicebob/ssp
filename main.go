package main

import (
	"log"
	"net/http"

	"github.com/alicebob/ssp/ssp"
)

const (
	listen = ":9998"
)

func main() {
	dsps := []ssp.DSP{
		{
			ID:     "dsp1",
			Name:   "Test 1",
			BidURL: "http://localhost:999/...",
		},
		{
			ID:     "dsp2",
			Name:   "Test 2",
			BidURL: "http://localhost:9999/...",
		},
	}
	pl1 := ssp.Placement{
		ID:   "my_website_1",
		Name: "My Website",
		// Type:   ssp.PlainBanner,
		// Width:  400,
		// Height: 200,
	}

	s := NewDaemon(dsps)
	// defer s.Close()
	// s.URLPrefix = "https://ssp.example/ssp"

	log.Printf("listening on %s...", listen)
	log.Fatal(http.ListenAndServe(listen, mux(s, []ssp.Placement{pl1})))

	// get(s.URL, "/p/my_website_1/code")
	// get(s.URL, "/p/my_website_1/example")
	// get(s.URL, "/p/my_website_1/embed.js")
	// get(s.URL, "/e/AAAA/view")
	// get(s.URL, "/e/AAAA/click")
}
