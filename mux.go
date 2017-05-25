package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alicebob/ssp/ssp"
	"github.com/julienschmidt/httprouter"
)

func mux(d *Daemon, pls []ssp.Placement) *httprouter.Router {
	r := httprouter.New()
	for _, pl := range pls {
		r.GET("/p/"+pl.ID+"/code.html", makeCode(pl))
		// r.GET("/p/"+pl.ID+"/embed.js", makeEmbed(d, &pl))
		r.GET("/p/"+pl.ID+"/iframe.html", makeIframe(d, &pl))
	}
	return r
}

func makeCode(pl ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		b, err := pl.Code()
		if err != nil {
			log.Printf("code: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(b))
	}
}

/*
func makeEmbed(d *Daemon, pl *ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		auc := d.RunAuction(pl)
		if auc == nil {
			// log.Printf("auction: no result")
			// http.Error(w, http.StatusText(500), 500)
			w.Header().Set("Content-Type", "application/javascript")
			// TODO: print some message in the banner
			return
		}
		b, err := pl.Embed(auc)
		if err != nil {
			log.Printf("auction: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(b))
	}
}
*/

func makeIframe(d *Daemon, pl *ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		auc := d.RunAuction(pl)
		if auc == nil {
			log.Printf("auction: no result")
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "no bid")
			return
		}
		b, err := pl.Iframe(auc)
		if err != nil {
			log.Printf("auction: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(b))
	}
}
