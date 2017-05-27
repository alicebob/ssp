package main

import (
	"context"
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
	a.FloorCPM = pl.FloorCPM
	a.Width = pl.Width
	a.Height = pl.Height

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	won, err := ssp.RunAuction(ctx, d.DSPs, a)
	if err != nil {
		log.Printf("bid error: %v", err)
		return a
	}
	if won == nil {
		return nil
	}
	log.Printf("winning bid: %s: %f", won.SSPID, won.PriceCPM)
	a.PriceCPM = won.PriceCPM
	a.NotificationURL = won.NotificationURL
	a.AdMarkup = won.AdMarkup
	a.Won()
	return a
}
