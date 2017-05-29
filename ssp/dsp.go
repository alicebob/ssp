package ssp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alicebob/ssp/openrtb"
)

var (
	ErrNoBid = errors.New("no bid")
)

type DSP struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	BidURL string `json:"bid_url"`
}

func (d *DSP) Bid(ctx context.Context, a *Auction) ([]Bid, error) {
	rtb := openrtb.BidRequest{
		ID: a.ID,
		Impressions: []openrtb.Impression{
			{
				ID:          "1",
				Bidfloor:    a.FloorCPM,
				BidfloorCur: Currency,
				Banner: &openrtb.Banner{
					Width:  a.Width,
					Height: a.Height,
				},
			},
		},
		Device: openrtb.Device{
			UserAgent: a.UserAgent,
			IP:        a.IP,
		},
		User: openrtb.User{
			ID: a.UserID,
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
	var bid openrtb.BidResponse
	if err := json.NewDecoder(resp.Body).Decode(&bid); err != nil {
		return nil, err
	}
	var bids []Bid
	for _, s := range bid.Seatbids {
		for _, b := range s.Bids {
			bids = append(bids, Bid{
				SSPID:           d.ID,
				PriceCPM:        b.Price,
				AdMarkup:        b.AdMarkup,
				NotificationURL: b.NotificationURL,
			})
		}
	}
	return bids, nil
}
