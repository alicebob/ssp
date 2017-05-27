package ssp

import (
	"github.com/alicebob/ssp/dsplib"
)

// RunDSP returns a http server you need to Close() when done
func RunDSP(id, name string, campaigns ...dsplib.Campaign) (DSP, *dsplib.DSP) {
	dsp := DSP{
		ID:   id,
		Name: name,
	}
	o := dsplib.NewDSP("localhost:0", campaigns)
	dsp.BidURL = o.BidURL // s.URL + "/rtb"
	return dsp, o
}
