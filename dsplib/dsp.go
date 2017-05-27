// Basic test DSP
package dsplib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Mux(cs []Campaign) *httprouter.Router {
	r := httprouter.New()
	r.POST("/rtb", rtbHandler(cs))
	return r
}

func rtbHandler(cs []Campaign) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("RTB request")
		var req RTBBidRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("req: %s", err)
			http.Error(w, http.StatusText(400), 400)
			return
		}

		// log.Printf("RTB request: %#v", req)
		var bids []RTBBid
		for _, imp := range req.Impressions {
			bids = append(bids, makeBid(imp, cs)...)
		}
		if len(bids) == 0 {
			log.Printf("no bids")
			w.WriteHeader(204)
			fmt.Fprintf(w, "{}")
			return
		}
		res := RTBBidResponse{
			ID: req.ID,
			Seatbids: []RTBSeatbid{
				{
					Bids: bids,
				},
			},
		}
		pl, err := json.Marshal(res)
		if err != nil {
			log.Printf("req: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if n, err := w.Write(pl); err != nil || n != len(pl) {
			log.Printf("req: %s (n: %d/%d)", err, n, len(pl))
		}
	}
}

// place a bid on every campaign which matches.
func makeBid(imp RTBImpression, cs []Campaign) []RTBBid {
	var bids []RTBBid
	for _, c := range cs {
		switch {
		case imp.Banner != nil:
			b := imp.Banner
			if b.Width == c.Width && b.Height == c.Height {
				bids = append(
					bids,
					RTBBid{
						ImpressionID: imp.ID,
						Price:        c.BidCPM,
						AdMarkup: fmt.Sprintf(
							`<a href="%s"><img src="%s" style="width:%dpx; height=%dpx"></a>`,
							c.ClickURL,
							c.ImageURL,
							c.Width,
							c.Height,
						),
					},
				)
			}
		}
	}
	return bids
}
