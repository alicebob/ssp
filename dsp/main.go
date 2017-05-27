package main

import (
	"log"

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
	s := dsplib.NewDSP(listen, campaigns)
	defer s.Close()
	log.Printf("configured campaigns:")
	for _, c := range campaigns {
		log.Printf(" - %dx%d: %s ($%.2f)", c.Width, c.Height, c.ID, c.BidCPM)
	}
	log.Printf("BidURL: %s", s.BidURL)
	for {
	}
}
