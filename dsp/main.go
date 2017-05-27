package main

import (
	"log"
	"net/http"

	"github.com/alicebob/ssp/dsplib"
)

const listen = "localhost:9990"

var (
	campaigns = []dsplib.Campaign{
		dsplib.Campaign{
			ID:       "camp1",
			Width:    466,
			Height:   214,
			BidCPM:   0.43,
			ImageURL: "https://imgs.xkcd.com/comics/debugger.png",
			ClickURL: "https://xkcd.com/1163/",
		},
		dsplib.Campaign{
			ID:       "camp2",
			Width:    300,
			Height:   330,
			BidCPM:   0.12,
			ImageURL: "https://imgs.xkcd.com/comics/duty_calls.png",
			ClickURL: "https://xkcd.com/386/",
		},
	}
)

func main() {
	log.Printf("BidURL: http://%s/rtb", listen)
	log.Fatal(http.ListenAndServe(listen, dsplib.Mux(campaigns)))
}
