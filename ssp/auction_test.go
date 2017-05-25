package ssp

import (
	"strings"
	"testing"
)

func TestAuction(t *testing.T) {
	pl := Placement{
		ID: "pl1",
	}
	a := NewAuction()
	a.UserAgent = "mozilla 1.2.3"
	a.IP = "1.2.3.4"
	a.PlacementID = pl.ID
	a.ViewURL = "https://my.example/viewme.cgi"
	a.ClickURL = "https://my.example/clickme.cgi"
	a.WinURL = "https://my.example/winme.cgi"
	html, err := pl.Embed(a)
	if err != nil {
		t.Fatal(err)
	}
	if want := "viewme"; !strings.Contains(html, want) {
		t.Fatalf("%q not found", want)
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
