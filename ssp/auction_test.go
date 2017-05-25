package ssp

import (
	"testing"
)

func TestAuction(t *testing.T) {
	pl := Placement{
		ID:       "pl1",
		ImgURL:   "https://site.example/image.jpg",
		ClickURL: "https://site.example/clicked.cgi",
	}
	a := NewAuction()
	a.UserAgent = "mozilla 1.2.3"
	a.IP = "1.2.3.4"
	a.PlacementID = pl.ID
	a.ClickURL = pl.ClickURL
	html := pl.HTML(a.ID)
	if html == "" {
		t.Fatalf("empty html")
	}

	// here be bidding
	// a.PriceCPM = 0.1 // TODO: millis
	if err := a.Win(); err != nil {
		t.Fatal(err)
	}

	if err := a.View(); err != nil {
		t.Fatal(err)
	}

	if err := a.Click(); err != nil {
		t.Fatal(err)
	}
}
