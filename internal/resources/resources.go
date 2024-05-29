package resources

import "embed"

//go:embed public
var Public embed.FS

//go:embed index.go.html
var IndexTemplate string
