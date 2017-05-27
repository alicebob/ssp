package ssp

type Auction struct {
	ID              string
	PlacementID     string
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
