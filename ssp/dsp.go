package ssp

import (
	"context"
	"errors"
)

type DSP struct {
	ID     string
	Name   string
	BidURL string
}

func (d *DSP) Bid(ctx context.Context, a *Auction) (*Bid, error) {
	return nil, errors.New("unimplemented")
}
