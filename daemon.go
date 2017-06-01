package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/alicebob/ssp/ssp"
)

const (
	timeout    = 100 * time.Millisecond
	cookieName = "uid"
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
func (d *Daemon) RunAuction(pl *ssp.Placement, r *http.Request, userID string) *ssp.Auction {
	a := ssp.NewAuction()
	a.UserAgent = r.UserAgent()
	if addr := r.RemoteAddr; addr != "" {
		a.IP, _, _ = net.SplitHostPort(addr)
	}
	a.UserID = userID
	a.PlacementID = pl.ID
	a.PlacementType = pl.Type
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

func getUserID(r *http.Request) string {
	if c, err := r.Cookie(cookieName); err == nil && c != nil {
		return c.Value
	}
	return ssp.RandomID(10)
}

func userID(w http.ResponseWriter, r *http.Request) string {
	userID := getUserID(r)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    userID,
		Path:     "/",
		MaxAge:   100 * 24 * 60 * 60,
		HttpOnly: true,
	})
	return userID
}
