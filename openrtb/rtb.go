package openrtb

// minimal OpenRTB 2.3

type BidRequest struct {
	ID          string       `json:"id"`
	Impressions []Impression `json:"imp,omitempty"`
	Device      Device       `json:"device"`
}

type Impression struct {
	ID          string  `json:"id,omitempty"`
	Banner      *Banner `json:"banner,omitempty"`
	Bidfloor    float64 `json:"bidfloor,omitempty"`
	BidfloorCur string  `json:"bidfloorcur,omitempty"`
}

type Banner struct {
	Width  int `json:"w,omitempty"`
	Height int `json:"h,omitempty"`
}

type Device struct {
	UserAgent string `json:"ua,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type BidResponse struct {
	ID       string    `json:"id,omitempty"`
	Seatbids []Seatbid `json:"seatbid,omitempty"`
}

type Seatbid struct {
	Bids []Bid `json:"bid,omitempty"`
}

type Bid struct {
	ImpressionID    string  `json:"impid,omitempty"`
	Price           float64 `json:"price,omitempty"`
	AdMarkup        string  `json:"adm,omitempty"`
	NotificationURL string  `json:"nurl,omitempty"`
}
