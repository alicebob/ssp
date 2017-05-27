package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/alicebob/ssp/ssp"
)

const (
	timeout = 100 * time.Millisecond
)

type Daemon struct {
	BaseURL string
	DSPs    []ssp.DSP
}

func NewDaemon(base string, dsps []ssp.DSP) *Daemon {
	return &Daemon{
		BaseURL: base,
		DSPs:    dsps,
	}
}

// RunAuction for a placement. Can take up to $timeout to run.
func (d *Daemon) RunAuction(pl *ssp.Placement) *ssp.Auction {
	a := ssp.NewAuction()
	a.UserAgent = "chromium 4.5.6"
	a.IP = "5.6.7.8"
	a.PlacementID = pl.ID
	a.Width = pl.Width
	a.Height = pl.Height

	// TODO: store auction

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	won, err := d.runBids(ctx, a, pl)
	if err != nil {
		log.Printf("bid error: %v", err)
		// TODO: update auction
		return a
	}
	log.Printf("winning bid: %s: %f", won.SSPID, won.PriceCPM)
	// TODO: update auction
	a.PriceCPM = won.PriceCPM
	a.NotificationURL = won.NotificationURL
	a.AdMarkup = won.AdMarkup
	a.Won()
	return a
}

func (d *Daemon) runBids(ctx context.Context, a *ssp.Auction, pl *ssp.Placement) (*ssp.Bid, error) {
	bids := make(chan []ssp.Bid, len(d.DSPs))
	open := 0
	for _, dsp := range d.DSPs {
		open++
		go func(dsp ssp.DSP) {
			bs, err := dsp.Bid(ctx, a)
			if err != nil {
				log.Printf("dsp %s bid err: %s", dsp.ID, err)
			}
			bids <- bs
		}(dsp)
	}

	var won ssp.Bid
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case bs := <-bids:
			for _, b := range bs {
				// TODO: match impressions
				if b.PriceCPM > won.PriceCPM {
					won = b
				}
			}
			open--
			if open == 0 {
				break loop
			}
		}
	}
	if won.PriceCPM == 0.0 {
		// TODO: what now?
		return nil, errors.New("no bid")
	}
	return &won, nil
}
