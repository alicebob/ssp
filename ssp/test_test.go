// Basic test DSP
package ssp

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
)

// RunDSP returns a http server you need to close when done
func RunDSP(id, name string) (*DSP, *httptest.Server) {
	dsp := &DSP{
		ID:   id,
		Name: name,
	}
	r := httprouter.New()
	r.POST("/rtb", rtbHandler(id))
	s := httptest.NewServer(r)
	dsp.BidURL = s.URL + "/rtb"
	return dsp, s
}

func rtbHandler(id string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("RTB request for %s", id)
		w.WriteHeader(204)
		fmt.Fprintf(w, "{}")
	}
}
