package dsplib

type Campaign struct {
	ID            string
	Width, Height int
	ImageURL      string
	ClickURL      string
	BidCPM        float64
}
