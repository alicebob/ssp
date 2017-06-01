package ssp

import (
	"bytes"
	"html/template"
)

type Type string

const (
	Banner Type = "banner"
	Video  Type = "video"
)

type Placement struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	FloorCPM float64 `json:"floor_cpm"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	Type     Type    `json:"type"`
}

func (p Placement) Code(base string) (string, error) {
	b := &bytes.Buffer{}
	args := struct {
		Base      string
		Static    string
		Placement Placement
	}{
		Base:      base,
		Static:    base + "../../static/",
		Placement: p,
	}
	switch p.Type {
	case Video:
		if err := videoCode.Execute(b, args); err != nil {
			return "", err
		}
	default:
		if err := plainCode.Execute(b, args); err != nil {
			return "", err
		}
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

	videoCode = template.Must(template.New("code").Parse(`
<script src="{{.Static}}html5vast.js"></script>
<link rel="stylesheet" href="{{.Static}}html5vast.css" type="text/css"/>
<video id="example" width="{{.Placement.Width}}" height="{{.Placement.Height}}" controls>
<source src="http://sample-videos.com/video/mp4/240/big_buck_bunny_240p_1mb.mp4" type="video/mp4" />
</video>
<script>
html5vast("example","{{.Base}}vast.xml",{
	ad_caption: 'Watch me! Click me!',
	media_bitrate_min: 0,
	media_bitrate_max: 999999
});
</script>
`))
)
