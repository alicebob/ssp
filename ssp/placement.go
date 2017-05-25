package ssp

// basic banner for now
type Placement struct {
	ID   string
	Name string
	// Width    int
	// Height   int
}

func (p Placement) Embed() string {
	return `
<div id="banner123" style="width:520px;height:100px"></div>
<script src="{{.Base}}embed.js"></script>
<script>document.getElementById("banner123").append(sspRun("{{.Base}}"))</script>
`
}

func (p Placement) HTML(auctionID string) string {
	return `
function sspRun(base) {
	var wrap = document.createElement("div");
	wrap.style.width="520px";
	wrap.style.height="100px";
	var view = document.createElement("img");
	view.style.width="0";
	view.style.height="0";
	wrap.appendChild(view);
	var img = document.createElement("img");
	img.style.width = "520px";
	img.style.height = "100px";
	img.src = "https://imgs.xkcd.com/s/a899e84.jpg";
	wrap.appendChild(img);
	return wrap;
}
`
}
