package ssp

import (
	"context"
	"testing"
)

func TestDSPNoBid(t *testing.T) {
	dsp, s := RunDSP("dsp", "My Second DSP")
	defer s.Close()

	a := NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = "myplacement"
	_, err := dsp.Bid(context.Background(), a)
	if have, want := err, ErrNoBid; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestDSP(t *testing.T) {
	dsp, s := RunDSP("dsp", "My Second DSP")
	defer s.Close()

	a := NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = "myplacement"
	bids, err := dsp.Bid(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := len(bids), 1; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	bid := bids[0]
	if have, want := bid.PriceCPM, 0.42; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
