package ssp

import (
	"context"
	"strings"
	"testing"

	"github.com/alicebob/ssp/dsplib"
)

func TestDSPNoBid(t *testing.T) {
	dsp, s := RunDSP("dsp", "My Second DSP", dsplib.Campaign{
		Width: 400, Height: 500,
	})
	defer s.Close()

	a := NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = "myplacement"
	a.PlacementType = Banner
	a.Width = 400
	a.Height = 123
	_, err := dsp.Bid(context.Background(), a)
	if have, want := err, ErrNoBid; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestDSP(t *testing.T) {
	dsp, s := RunDSP("dsp", "My Second DSP", dsplib.Campaign{
		Width: 400, Height: 500, BidCPM: 0.42,
	})
	defer s.Close()

	a := NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = "myplacement"
	a.PlacementType = Banner
	a.Width = 400
	a.Height = 500
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

func TestDSPMultiple(t *testing.T) {
	dsp, s := RunDSP("dsp", "My Second DSP",
		dsplib.Campaign{
			Width: 400, Height: 500, BidCPM: 0.42,
		},
		dsplib.Campaign{
			Width: 400, Height: 500, BidCPM: 0.52,
		},
		dsplib.Campaign{
			Width: 400, Height: 500, BidCPM: 0.22,
		},
	)
	defer s.Close()

	a := NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = "myplacement"
	a.PlacementType = Banner
	a.Width = 400
	a.Height = 500
	bids, err := dsp.Bid(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := len(bids), 3; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}
	bid := bids[0]
	if have, want := bid.PriceCPM, 0.42; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestVideo(t *testing.T) {
	dsp, s := RunDSP("dsp", "My Second DSP", dsplib.Campaign{
		Width: 400, Height: 500, BidCPM: 0.42,
		Type:     "video",
		VideoURL: "http://some.where/movie.mp4",
	})
	defer s.Close()

	a := NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = "myplacement"
	a.PlacementType = Video
	a.Width = 400
	a.Height = 500
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
	if want, have := "movie.mp4", bid.AdMarkup; !strings.Contains(have, want) {
		t.Fatalf("%q not found", want)
	}
}
