package netdriver

import (
	"html/template"

	"github.com/fluffle/sp0rkle/util/push"
)

type pushDevice struct {
	State   string
	Devices []*push.Device
}

var pushDeviceTmpl = template.Must(template.New("pushdevice").Parse(`<html>
<head>
  <title>sp0rkle's shonky device chooser</title>
</head>
<body>
  <h1>Choose a device to send sp0rkle pushes to from the list below.</h1>
  <form action="/oauth/device" method="POST">
    <input type="hidden" name="state" value="{{.State}}">
{{ range $i, $d := .Devices }}
    <input type="radio" name="iden" value="{{ $d.Iden }}" {{ if not $i }}checked{{ end }}>
    {{ $d.Nickname }} -- {{ $d.Manufacturer }} {{ $d.Model }}
    ({{ if not $d.Active }}in{{ end }}active,
    {{ if not $d.Pushable }}not {{ end }}pushable)<br />
{{ end }}
    <input type="submit" value="Choose Device">
  </form>
</body>
</html>`))

var pushSuccessHTML = `<html>
<head>
  <title>You are teh winnar!</title>
</head>
<body>
  <h1>YAY.</h1>
  <p>sp0rkle has successfully negotiated the OAuth dance, and you ought
  to be receiving a confirmation push any time soon. Simply tell sp0rkle
  "push auth &lt;pin&gt;" to complete the setup.</p>
</body>
</html>`

var pushFailures = map[string]string{
	"parse":    "Parsing HTTP form values failed.",
	"denied":   "Authorization denied by user.",
	"exchange": "Could not exchange response for token.",
	"device":   "Could not get list of devices.",
	"push":     "Could not push confirmation message.",
	"notoken":  "No authorization token found.",
	"nostate":  "Bad or missing oauth state.",
	"noiden":   "Bad or missing device iden.",
	"nofail":   "No failure reason provided, suckah.",
}

type pushFailure struct {
	Message string
}

var pushFailureTmpl = template.Must(template.New("pushfailure").Parse(`<html>
<head>
  <title>You are teh losar!</title>
</head>
<body>
  <h1>Oh noes!</h1>
  <p>It's all gone wrong. You'd better moan at fluffle.</p>
  <p>{{ .Message }}</p>
</body>
</html>`))
