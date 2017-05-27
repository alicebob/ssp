package ssp

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Auction struct {
	ID              string
	PlacementID     string
	FloorCPM        float64
	Width, Height   int
	UserAgent       string
	IP              string
	PriceCPM        float64
	AdMarkup        string
	NotificationURL string
}

func NewAuction() *Auction {
	return &Auction{
		ID: "123", // TODO
	}
}

// Won will call the notication callback, if any
func (a *Auction) Won() {
	n := a.NotificationURL
	if n == "" {
		return
	}
	url := strings.Replace(n, "${AUCTION_PRICE}", fmt.Sprintf("%0.2f", a.PriceCPM), -1)
	go func() {
		// TODO: retry?
		res, err := http.Get(url)
		if err != nil {
			log.Printf("notification err: %s", err)
			return
		}
		defer res.Body.Close()
	}()
}
