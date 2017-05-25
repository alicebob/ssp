package main

import (
	"log"

	"github.com/alicebob/ssp/ssp"
)

func main() {
	dsp, s := ssp.RunDSP("dsp1", "My DSP")
	defer s.Close()
	log.Printf("BidURL: %s", dsp.BidURL)
	for {
	}
}
