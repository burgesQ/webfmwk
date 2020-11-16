package docs

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/burgesQ/webfmwk/v4"
)

// TODO
type RedocParam struct {
	Path   string
	DocURI string
}

var (
	_redocTmpl = `<!DOCTYPE html>
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

	_defRedoc = &RedocParam{
		DocURI: "/api/docs/swagger.json",
		Path:   "/docs/redoc",
	}
)

func GetTemplate() string {
	return _redocTmpl
}

func SetTemplate(tmpl string) {
	_redocTmpl = tmpl
}

func (p *RedocParam) sync() {
	if p.DocURI == "" {
		p.DocURI = "/api/docs/swagger.json"

	}

	if p.Path == "" {
		p.Path = "/docs/redoc"
	}
}

// Return a DocHandler settup for redoc
// use of template, params expect the DocURI string
func GetRedocHandler(p *RedocParam) webfmwk.DocHandler {
	if p == nil {
		p = _defRedoc
	}
	p.sync()

	t := template.Must(template.New("redoc").Parse(_redocTmpl))
	buf := bytes.NewBuffer(nil)
	_ = t.Execute(buf, p)
	b := buf.Bytes()

	return webfmwk.DocHandler{
		Name: "redoc",
		Path: p.Path,
		H: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(b)
		},
	}
}
