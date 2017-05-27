package ssp

import (
	"testing"
)

func TestPickBidSimple(t *testing.T) {
	bs := []Bid{
		{
			SSPID:    "ssp1",
			PriceCPM: 0.4,
		},
		{
			SSPID:    "ssp2",
			PriceCPM: 0.6, // 2nd price
		},
		{
			SSPID:    "ssp3",
			PriceCPM: 0.8, // Winner
		},
	}
	if have, want := *pickBid(bs), (Bid{
		SSPID:    "ssp3",
		PriceCPM: 0.6,
	}); have != want {
		t.Errorf("have %+v, want %+v", have, want)
	}
}

func TestPickBidNoBid(t *testing.T) {
	bs := []Bid{}
	if have, want := pickBid(bs), (*Bid)(nil); have != want {
		t.Errorf("have %+v, want %+v", have, want)
	}
}

func TestPickBidSingle(t *testing.T) {
	bs := []Bid{
		{
			SSPID:    "ssp1",
			PriceCPM: 0.4,
		},
	}
	if have, want := *pickBid(bs), (Bid{
		SSPID:    "ssp1",
		PriceCPM: 0.0, // no second bid
	}); have != want {
		t.Errorf("have %+v, want %+v", have, want)
	}
}
