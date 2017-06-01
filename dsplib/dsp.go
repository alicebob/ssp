// Basic test DSP
package dsplib

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/ssp/openrtb"
)

type DSP struct {
	BaseURL   string
	BidURL    string
	server    *http.Server
	campaigns []Campaign
	wonCount  int
	wonCPM    float64
	mu        sync.Mutex
}

func NewDSP(listen string, cs []Campaign) *DSP {
	l, err := net.Listen("tcp", listen)
	if err != nil {
		panic(err.Error())
	}
	port := l.Addr().(*net.TCPAddr).Port
	base := fmt.Sprintf("http://localhost:%d/", port)
	d := &DSP{
		BaseURL:   base,
		BidURL:    fmt.Sprintf("%srtb", base),
		campaigns: cs,
	}
	d.server = &http.Server{
		Addr:    listen,
		Handler: d.Mux(),
	}
	go d.server.Serve(l)
	return d
}

func (dsp *DSP) Close() error {
	return dsp.server.Close()
}

func (dsp *DSP) Mux() *httprouter.Router {
	r := httprouter.New()
	r.POST("/rtb", dsp.rtbHandler())
	r.GET("/win", dsp.winHandler())
	return r
}

func (dsp *DSP) rtbHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("RTB request")
		var req openrtb.BidRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("req: %s", err)
			http.Error(w, http.StatusText(400), 400)
			return
		}

		// log.Printf("RTB request: %#v", req)
		var bids []openrtb.Bid
		for _, imp := range req.Impressions {
			bids = append(bids, dsp.makeBid(imp)...)
		}
		if len(bids) == 0 {
			log.Printf("no bids")
			w.WriteHeader(204)
			fmt.Fprintf(w, "{}")
			return
		}
		res := openrtb.BidResponse{
			ID: req.ID,
			Seatbids: []openrtb.Seatbid{
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
func (dsp *DSP) makeBid(imp openrtb.Impression) []openrtb.Bid {
	var bids []openrtb.Bid
	for _, c := range dsp.campaigns {
		switch {
		case imp.Banner != nil && (c.Type == "banner" || c.Type == ""):
			b := imp.Banner
			if b.Width == c.Width && b.Height == c.Height {
				bids = append(
					bids,
					openrtb.Bid{
						ImpressionID: imp.ID,
						Price:        c.BidCPM,
						AdMarkup: fmt.Sprintf(
							`<a href="%s"><img src="%s" style="width:%dpx; height=%dpx"></a>`,
							c.ClickURL,
							c.ImageURL,
							c.Width,
							c.Height,
						),
						NotificationURL: dsp.winURL(),
					},
				)
			}
		case imp.Video != nil && c.Type == "video":
			// TODO: fuzzier dimension matching
			v := imp.Video
			if v.Width == c.Width && v.Height == c.Height {
				bids = append(
					bids,
					openrtb.Bid{
						ImpressionID:    imp.ID,
						Price:           c.BidCPM,
						AdMarkup:        vast30(c),
						NotificationURL: dsp.winURL(),
					},
				)
			}
		}
	}
	return bids
}

func (dsp *DSP) winHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		pr := r.FormValue("p")
		if pr == "" {
			log.Printf("no price")
			http.Error(w, http.StatusText(400), 400)
			return
		}
		price, err := strconv.ParseFloat(pr, 64)
		if err != nil {
			log.Printf("invalid price %q: %s", pr, err)
			http.Error(w, http.StatusText(400), 400)
			return
		}
		log.Printf("won %f", price)
		dsp.mu.Lock()
		defer dsp.mu.Unlock()
		dsp.wonCount++
		dsp.wonCPM += price
	}
}

func (dsp *DSP) winURL() string {
	return fmt.Sprintf("%swin?p=${AUCTION_PRICE}", dsp.BaseURL)
}

// Won returns count+total CPM of winnotices
func (dsp *DSP) Won() (int, float64) {
	dsp.mu.Lock()
	defer dsp.mu.Unlock()
	return dsp.wonCount, dsp.wonCPM
}

func vast30(c Campaign) string {
	return fmt.Sprintf(`<VAST version="3.0"><Ad id="%s"><InLine><AdSystem>My First SSP</AdSystem><AdTitle>%s</AdTitle><Creatives><Creative><Linear><Duration>00:01:00.000</Duration><MediaFiles><MediaFile delivery="progressive" width="%d" height="%d" type="video/mp4" bitrate="1000"><![CDATA[%s]]></MediaFile></MediaFiles><VideoClicks><ClickThrough>%s</ClickThrough></VideoClicks></Linear></Creative></Creatives></InLine></Ad></VAST>`,
		c.ID,
		c.ID,
		c.Width,
		c.Height,
		c.VideoURL,
		c.ClickURL,
	)
}
