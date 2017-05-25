package ssp

// TODO: lock
type Auction struct {
	ID          string
	PlacementID string
	UserAgent   string
	IP          string
	ClickURL    string
	WinURL      string
	Wins        int
	Views       int
	Clicks      int
}

func NewAuction() *Auction {
	return &Auction{
		ID: "123",
	}
}

func (a *Auction) Win() error {
	// TODO: call WinURL
	a.Wins++
	return nil
}

func (a *Auction) View() error {
	a.Views++
	return nil
}

func (a *Auction) Click() error {
	a.Clicks++
	return nil
}
