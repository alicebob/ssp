package ssp

type Bid struct {
	SSPID           string
	PriceCPM        float64 // TODO: millis
	AdMarkup        string  // snippet
	NotificationURL string
}
