package schemas

import "embed"

//go:embed *.schema.json
var FileFS embed.FS
