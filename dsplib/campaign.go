package dsplib

type Campaign struct {
	ID            string
	Type          string
	Width, Height int
	ImageURL      string
	VideoURL      string
	ClickURL      string
	BidCPM        float64
}
