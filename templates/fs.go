package templates

import "embed"

//go:embed *.tmpl galleries/*.tmpl
var FS embed.FS
