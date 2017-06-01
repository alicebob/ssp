package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alicebob/ssp/dsplib"
	"github.com/alicebob/ssp/openrtb"
	"github.com/alicebob/ssp/ssp"
)

var (
	camp1 = dsplib.Campaign{
		ID:       "camp1",
		Type:     "banner",
		Width:    466,
		Height:   214,
		BidCPM:   0.43,
		ImageURL: "https://imgs.xkcd.com/comics/debugger.png",
		ClickURL: "https://xkcd.com/1163/",
	}
	camp2 = dsplib.Campaign{
		ID:       "camp2",
		Type:     "banner",
		Width:    300,
		Height:   330,
		BidCPM:   0.12,
		ImageURL: "https://imgs.xkcd.com/comics/duty_calls.png",
		ClickURL: "https://xkcd.com/386/",
	}
	pl1 = ssp.Placement{
		ID:       "my_website_1",
		Type:     ssp.Banner,
		Name:     "My Website",
		FloorCPM: 0.2,
		Width:    466,
		Height:   214,
	}
)

func TestMain(t *testing.T) {
	jar, _ := cookiejar.New(nil)
	cl := &http.Client{
		Jar: jar,
	}

	dsp1, s1 := ssp.RunDSP("dsp1", "My First DSP")
	defer s1.Close()
	dsp2, s2 := ssp.RunDSP("dsp2", "My Second DSP", camp1, camp2)
	defer s2.Close()

	d := NewDaemon("http://localhost/", []ssp.DSP{dsp1, dsp2})
	s := httptest.NewServer(mux(d, []ssp.Placement{pl1}))
	defer s.Close()

	{
		r := getok(t, s, cl, 200, "/")
		if want := "My Website"; !strings.Contains(r, want) {
			t.Errorf("not found: %q", want)
		}
	}

	{
		r := getok(t, s, cl, 200, "/p/my_website_1/code.html")
		if want := "<iframe"; !strings.Contains(r, want) {
			t.Errorf("not found: %q", want)
		}
	}

	{
		r := getok(t, s, cl, 200, "/p/my_website_1/iframe.html")
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

	getok(t, s, cl, 404, "/p/my_website_1/foo.html")
	getok(t, s, cl, 404, "/p/foo/code.html")
}

func TestRTB(t *testing.T) {
	jar, _ := cookiejar.New(nil)
	cl := &http.Client{
		Jar: jar,
	}

	var lastReq openrtb.BidRequest
	r := httprouter.New()
	r.POST("/rtb", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if err := json.NewDecoder(r.Body).Decode(&lastReq); err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		w.WriteHeader(204)
		fmt.Fprintf(w, "{}")
	})
	dspS := httptest.NewServer(r)
	defer dspS.Close()
	dsp := ssp.DSP{
		ID:     "dsp1",
		Name:   "Hello world",
		BidURL: dspS.URL + "/rtb",
	}

	d := NewDaemon("http://localhost/", []ssp.DSP{dsp})
	s := httptest.NewServer(mux(d, []ssp.Placement{pl1}))
	defer s.Close()

	getok(t, s, cl, 200, "/p/my_website_1/iframe.html")
	if have, want := len(lastReq.Impressions), 1; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	if have, want := lastReq.Device.UserAgent, "Go-http-client/1.1"; have != want {
		t.Fatalf("have %s, want %s", have, want)
	}
	if have, want := lastReq.Device.IP, "127.0.0.1"; have != want {
		t.Fatalf("have %s, want %s", have, want)
	}
	userID := lastReq.User.ID
	if have := userID; have == "" {
		t.Fatalf("empty value")
	}

	// userid should be stable
	getok(t, s, cl, 200, "/p/my_website_1/iframe.html")
	if have, want := lastReq.User.ID, userID; have != want {
		t.Fatalf("have %s, want %s", have, want)
	}
}

func getok(t *testing.T, s *httptest.Server, cl *http.Client, status int, path string) string {
	req, err := http.NewRequest("GET", s.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := cl.Do(req)
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
