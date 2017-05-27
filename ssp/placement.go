package ssp

import (
	"bytes"
	"html/template"
)

// basic banner
type Placement struct {
	ID       string
	Name     string
	FloorCPM float64
	Width    int
	Height   int
}

func (p Placement) Code(base string) (string, error) {
	b := &bytes.Buffer{}
	args := struct {
		Base      string
		Placement Placement
	}{
		Base:      base,
		Placement: p,
	}
	if err := plainCode.Execute(b, args); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (p Placement) Iframe(au *Auction) (string, error) {
	b := &bytes.Buffer{}
	markup := template.HTML(au.AdMarkup)
	if err := plainIframe.Execute(b, markup); err != nil {
		return "", err
	}
	return b.String(), nil
}

var (
	plainCode = template.Must(template.New("code").Parse(`
<iframe src="{{.Base}}iframe.html" style="border: 0; width: {{.Placement.Width}}px; height: {{.Placement.Height}}px"></iframe>
`))

	plainIframe = template.Must(template.New("embed").Parse(`
<html>
<body style="margin: 0">
{{.}}
<script>
window.onload = function() {
  var as = document.getElementsByTagName('a');
  for (var i=0; i<as.length; i++){
    as[i].setAttribute('target', '_top');
  }
}
</script>
`))
)
