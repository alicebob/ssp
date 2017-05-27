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
	a.AdMarkup = `<img src="http://my.example/image.png">`
	a.NotificationURL = "https://my.example/winme.cgi"
	html, err := pl.Iframe(a)
	if err != nil {
		t.Fatal(err)
	}
	if want := "image.png"; !strings.Contains(html, want) {
		t.Fatalf("%q not found", want)
	}
}
