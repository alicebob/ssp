package ssp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNoBid = errors.New("no bid")
)

type DSP struct {
	ID     string
	Name   string
	BidURL string
}

func (d *DSP) Bid(ctx context.Context, a *Auction) ([]Bid, error) {
	rtb := RTBBidRequest{
		ID: a.ID,
		Impressions: []RTBImpression{
			{
				ID: "1",
				Banner: &RTBBanner{
					Width:  a.Width,
					Height: a.Height,
				},
			},
		},
	}
	pl, err := json.Marshal(rtb)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{}
	cl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest("POST", d.BidURL, bytes.NewBuffer(pl))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	switch s := resp.StatusCode; s {
	case 204:
		return nil, ErrNoBid
	case 200:
		// pass
	default:
		return nil, fmt.Errorf("unexpected HTTP status code: %d %s", s, resp.Status)
	}
	var bid RTBBidResponse
	if err := json.NewDecoder(resp.Body).Decode(&bid); err != nil {
		return nil, err
	}
	b := bid.Seatbids[0].Bids[0] // TODO
	return []Bid{
		{
			SSPID:           d.ID,
			PriceCPM:        b.Price,
			AdMarkup:        b.AdMarkup,
			NotificationURL: b.NotificationURL,
		},
	}, nil
}
