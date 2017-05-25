// Basic test DSP
package ssp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
)

// RunDSP returns a http server you need to Close() when done
func RunDSP(id, name string) (*DSP, *httptest.Server) {
	dsp := &DSP{
		ID:   id,
		Name: name,
	}
	s := httptest.NewServer(Mux())
	dsp.BidURL = s.URL + "/rtb"
	return dsp, s
}

func Mux() *httprouter.Router {
	r := httprouter.New()
	r.POST("/rtb", rtbHandler())
	return r
}

func rtbHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("RTB request")
		var req RTBBidRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("req: %s", err)
			http.Error(w, http.StatusText(400), 400)
			return
		}

		// log.Printf("RTB request for %s: %#v", id, req)
		if len(req.Impressions) == 0 {
			log.Printf("no impressions")
			w.WriteHeader(204)
			fmt.Fprintf(w, "{}")
			return
		}
		imp := req.Impressions[0]

		res := RTBBidResponse{
			ID:    req.ID,
			BidID: "456", // TODO
			Seatbids: []RTBSeatbid{
				{
					Bids: []RTBBid{
						{
							ID:           "1",
							ImpressionID: imp.ID,
							Price:        0.42,
							AdMarkup: `
<a href="https://xkcd.com?clicked"><img src="https://imgs.xkcd.com/s/a899e84.jpg"></a>
`,
						},
					},
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
