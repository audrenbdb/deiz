package mailtmpl

import (
	"embed"
	"text/template"
)

//go:embed *.html
var emailTemplates embed.FS

func Embed() (*template.Template, error) {
	return template.ParseFS(emailTemplates, "*")
}
