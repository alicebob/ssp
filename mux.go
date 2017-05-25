package main

import (
	"fmt"
	"net/http"

	"github.com/alicebob/ssp/ssp"
	"github.com/julienschmidt/httprouter"
)

func mux(d *Daemon, pls []ssp.Placement) *httprouter.Router {
	r := httprouter.New()
	for _, pl := range pls {
		r.GET("/p/"+pl.ID+"/code.html", makeCode(pl))
		r.GET("/p/"+pl.ID+"/embed.js", makeEmbed(d, &pl))
	}
	return r
}

func makeCode(pl ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, pl.Code())
	}
}

func makeEmbed(d *Daemon, pl *ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/javascript")
		auc := d.RunAuction(pl)
		if auc == nil {
			// log.Printf("auction: no result")
			// http.Error(w, http.StatusText(500), 500)
			return
		}
		fmt.Fprintf(w, pl.Embed(auc.ID))
	}
}
