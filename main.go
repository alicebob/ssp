package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alicebob/ssp/ssp"
)

const (
	listen = "localhost:9998"
)

type Config struct {
	DSPs       []ssp.DSP       `json:"dsps"`
	Placements []ssp.Placement `json:"placements"`
}

var (
	config = flag.String("config", "./ssp.json", "config file")
)

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatal(err)
	}
	var c Config
	if err := json.Unmarshal(f, &c); err != nil {
		log.Fatal(err)
	}
	log.Printf("configured DSPs:")
	for _, d := range c.DSPs {
		log.Printf(" - %s", d.Name)
	}
	s := NewDaemon(fmt.Sprintf("http://%s", listen), c.DSPs)
	log.Printf("listening on %s...", listen)
	log.Fatal(http.ListenAndServe(listen, mux(s, c.Placements)))
}
