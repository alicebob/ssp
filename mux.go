package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/alicebob/ssp/ssp"
	"github.com/julienschmidt/httprouter"
)

func mux(d *Daemon, pls []ssp.Placement) *httprouter.Router {
	r := httprouter.New()
	r.GET("/", makeList(pls))
	for _, pl := range pls {
		base := "/p/" + pl.ID + "/"
		r.GET(base, makeExample(d.BaseURL+base, pl))
		r.GET(base+"code.html", makeCode(d.BaseURL+base, pl))
		r.GET(base+"iframe.html", makeIframe(d, pl))
	}
	return r
}

func makeList(pls []ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "text/html")
		runTemplate(w, listTemplate, pls)
	}
}

func makeExample(base string, pl ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		b, err := pl.Code(base)
		if err != nil {
			log.Printf("code: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		args := map[string]interface{}{
			"id":     pl.ID,
			"name":   pl.Name,
			"width":  pl.Width,
			"height": pl.Height,
			"code":   template.HTML(string(b)),
		}
		w.Header().Set("Content-Type", "text/html")
		runTemplate(w, exampleTemplate, args)
	}
}

func makeCode(base string, pl ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		b, err := pl.Code(base)
		if err != nil {
			log.Printf("code: %s", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(b))
	}
}

func makeIframe(d *Daemon, pl ssp.Placement) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		userID := getUserID(r)
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    userID,
			Path:     "/",
			MaxAge:   100 * 24 * 60 * 60,
			HttpOnly: true,
		})

		auc := d.RunAuction(&pl, r, userID)
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

func runTemplate(w http.ResponseWriter, t *template.Template, args interface{}) {
	b := &bytes.Buffer{}
	if err := t.Execute(b, args); err != nil {
		log.Printf("template: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if n, err := w.Write(b.Bytes()); err != nil || n != b.Len() {
		log.Printf("write error (%d/%d): %s", n, b.Len(), err)
		return
	}
}

var (
	listTemplate = template.Must(template.New("list").Parse(`
<html>
<title>Placement list</title>
<body>
Available placements:<br />
<br />
{{range .}}
	<a href="/p/{{ .ID }}/">{{.Name}}</a><br />
{{end}}
`))

	exampleTemplate = template.Must(template.New("list").Parse(`
<html>
<title>{{ .name }}</title>
<body>
<b>{{ .name}}</b><br />
{{ .width }}x{{ .height }}<br />
<br />
<br />

Embed code:<br />
<code style="background-color: #eee">
	{{ .code | html }}
</code>
<br />
<a href="./code.html">raw</a>
<br />
<br />
<br />

<div style="width:{{.width}}px; height:{{.height}}px; border: solid 1px gray">
{{ .code }}
</div>
`))
)
