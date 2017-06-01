package ssp

import (
	"context"
	"log"
	"sort"
	"sync"
)

type Bid struct {
	SSPID           string
	PriceCPM        float64
	NotificationURL string
	Type            Type
	AdMarkup        string // VAST also goes here
}

func RunAuction(ctx context.Context, dsps []DSP, a *Auction) (*Bid, error) {
	var (
		wg   sync.WaitGroup
		bids = make(chan Bid)
	)
	for _, dsp := range dsps {
		wg.Add(1)
		go func(dsp DSP) {
			defer wg.Done()
			bs, err := dsp.Bid(ctx, a)
			if err != nil {
				log.Printf("dsp %s bid err: %s", dsp.ID, err)
			}
			for _, b := range bs {
				bids <- b
			}
		}(dsp)
	}
	go func() {
		wg.Wait()
		close(bids)
	}()

	var allBids []Bid
	for b := range bids {
		// TODO: match impression slots
		if b.PriceCPM < a.FloorCPM {
			continue
		}
		allBids = append(allBids, b)
	}
	won := pickBid(allBids)
	if won != nil && won.PriceCPM == 0 {
		won.PriceCPM = a.FloorCPM
	}
	return won, nil
}

// pickBid with a second price auction. Will return a 0 price if there is only
// one auction. Does not check validness of any bid.
func pickBid(bs []Bid) *Bid {
	if len(bs) == 0 {
		return nil
	}
	sort.Slice(bs, func(i, j int) bool {
		if bs[i].PriceCPM == bs[j].PriceCPM {
			return bs[i].SSPID > bs[j].SSPID
		}
		return bs[i].PriceCPM > bs[j].PriceCPM
	})
	won := bs[0]
	won.PriceCPM = 0.0
	if len(bs) > 1 {
		won.PriceCPM = bs[1].PriceCPM
	}
	return &won
}
