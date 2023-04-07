package redoc

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/burgesQ/webfmwk/v5"
)

const (
	_defPath = "api/doc/swagger.json"
	_defURI  = "/doc/redoc"
)

// Param hold the required metadata to expose the redoc handler.
type Param struct {
	// Path hold the value on which is exposed the handler
	Path string
	// DocURI hold the swagger.json URI
	DocURI string
}

var _defRedoc = &Param{
	DocURI: _defURI,
	Path:   _defPath,
}

// Path set the redoc handler path.
func Path(path string) func(*Param) {
	return func(rp *Param) {
		rp.Path = path
	}
}

// DocURI set the redoc handle source swagger url.
func DocURI(uri string) func(*Param) {
	return func(rp *Param) {
		rp.DocURI = uri
	}
}

// GetHandler return a DocHandler settup for redoc
// use of template, params expect the DocURI string.
//
//	var opts = []webfmwk.Option{
//		webfmwk.WithDocHandlers(redoc.GetHandler(
//			redoc.DocURI("/api/v2/docs/swagger.json")
//		)
//	)}
func GetHandler(opt ...func(*Param)) webfmwk.DocHandler {
	p := _defRedoc

	for _, o := range opt {
		o(p)
	}

	b := genContent(p)

	return webfmwk.DocHandler{
		Name: "redoc",
		Path: p.Path,
		H: func(c webfmwk.Context) error {
			c.SetContentType("text/html; charset=utf-8")
			c.SetStatusCode(http.StatusOK)
			c.GetFastContext().Response.SetBody(b)

			return nil
		},
	}
}

func genContent(p *Param) []byte {
	t := template.Must(template.New("redoc").Parse(_redocTmpl))
	buf := bytes.NewBuffer(nil)
	_ = t.Execute(buf, p)

	return buf.Bytes()
}

const (
	_redocTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<title>ReDoc</title>
		<!-- needed for adaptive design -->
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

		<!--
				ReDoc doesn't change outer page styles
			-->
		<style>
			body {
			margin: 0;
			padding: 0;
			}
		</style>
	</head>
	<body>
		<redoc spec-url={{ .DocURI }}></redoc>
		<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
	</body>
</html>`
)
