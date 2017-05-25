package ssp

import (
	"bytes"
	"html/template"
)

// basic banner
type Placement struct {
	ID     string
	Name   string
	Width  int
	Height int
}

func (p Placement) Code() (string, error) {
	b := &bytes.Buffer{}
	if err := plainCode.Execute(b, p); err != nil {
		return "", err
	}
	return b.String(), nil
}

/*
func (p Placement) Embed(au *Auction) (string, error) {
	b := &bytes.Buffer{}
	markup := template.JS(oneLine(au.AdMarkup))
	if err := plainEmbed.Execute(b, au, markup); err != nil {
		return "", err
	}
	return b.String(), nil
}
*/

func (p Placement) Iframe(au *Auction) (string, error) {
	b := &bytes.Buffer{}
	markup := template.HTML(au.AdMarkup)
	if err := plainIframe.Execute(b, markup); err != nil {
		return "", err
	}
	return b.String(), nil
}

var plainCode = template.Must(template.New("code").Parse(`
<iframe width="{{.Width}}" height="{{.Height}}" src=".../iframe.html" style="border: 0"></iframe>
`))

/*
var plainEmbed = template.Must(template.New("embed").Parse(`
function sspRun() {
	console.log("rendering auction", "{{.AdMarkup}}");
	var d = document.createElement("div");
	d.style.width = "{{.Width}}px";
	d.style.height = "{{.Height}}px";
	d.innerHTML = "{{.AdMarkup}}";
	return d;
}
`))
*/

var plainIframe = template.Must(template.New("embed").Parse(`
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
