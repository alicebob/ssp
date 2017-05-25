package ssp

type Bid struct {
	SSPID    string
	PriceCPM float64 // TODO: millis
	ImageURL string
	ClickURL string
	// TODO: view url
	// TODO: win url
}
