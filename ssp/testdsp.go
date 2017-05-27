package ssp

import (
	"net/http/httptest"

	"github.com/alicebob/ssp/dsplib"
)

// RunDSP returns a http server you need to Close() when done
func RunDSP(id, name string, campaigns ...dsplib.Campaign) (DSP, *httptest.Server) {
	dsp := DSP{
		ID:   id,
		Name: name,
	}
	s := httptest.NewServer(dsplib.Mux(campaigns))
	dsp.BidURL = s.URL + "/rtb"
	return dsp, s
}
