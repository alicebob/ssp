package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/ssp/dsplib"
	"github.com/alicebob/ssp/ssp"
)

var (
	camp1 = dsplib.Campaign{
		ID:       "camp1",
		Width:    466,
		Height:   214,
		BidCPM:   0.43,
		ImageURL: "https://imgs.xkcd.com/comics/debugger.png",
		ClickURL: "https://xkcd.com/1163/",
	}
	camp2 = dsplib.Campaign{
		ID:       "camp2",
		Width:    300,
		Height:   330,
		BidCPM:   0.12,
		ImageURL: "https://imgs.xkcd.com/comics/duty_calls.png",
		ClickURL: "https://xkcd.com/386/",
	}
	pl1 = ssp.Placement{
		ID:       "my_website_1",
		Name:     "My Website",
		FloorCPM: 0.2,
		Width:    466,
		Height:   214,
	}
)

func TestMain(t *testing.T) {
	dsp1, s1 := ssp.RunDSP("dsp1", "My First DSP")
	defer s1.Close()
	dsp2, s2 := ssp.RunDSP("dsp2", "My Second DSP", camp1, camp2)
	defer s2.Close()

	d := NewDaemon("http://localhost/", []ssp.DSP{dsp1, dsp2})
	s := httptest.NewServer(mux(d, []ssp.Placement{pl1}))
	defer s.Close()

	{
		r := getok(t, s, 200, "/")
		if want := "My Website"; !strings.Contains(r, want) {
			t.Errorf("not found: %q", want)
		}
	}

	{
		r := getok(t, s, 200, "/p/my_website_1/code.html")
		if want := "<iframe"; !strings.Contains(r, want) {
			t.Errorf("not found: %q", want)
		}
	}

	{
		r := getok(t, s, 200, "/p/my_website_1/iframe.html")
		if have, want := r, "debugger.png"; !strings.Contains(have, want) {
			t.Errorf("not found: %q", want)
		}
		if have, want := r, "214px"; !strings.Contains(have, want) {
			t.Errorf("not found: %q in %q", want, r)
		}
		time.Sleep(10 * time.Millisecond)
		{
			wonCount, wonCMP := s2.Won()
			if have, want := wonCount, 1; have != want {
				t.Errorf("have %v, want %v", have, want)
			}
			// Fallback to floor
			if have, want := wonCMP, 0.2; have != want {
				t.Errorf("have %v, want %v", have, want)
			}
		}
	}

	getok(t, s, 404, "/p/my_website_1/foo.html")
	getok(t, s, 404, "/p/foo/code.html")
}

func getok(t *testing.T, s *httptest.Server, status int, path string) string {
	url := s.URL + path
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if have, want := res.StatusCode, status; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
