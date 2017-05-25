package ssp

import (
	"context"
	"testing"
)

func TestDSP(t *testing.T) {
	dsp, s := RunDSP("dsp", "My Second DSP")
	defer s.Close()

	a := NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = "myplacement"
	ctx := context.Background()
	bid, err := dsp.Bid(ctx, a)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := bid.PriceCPM, 0.42; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
