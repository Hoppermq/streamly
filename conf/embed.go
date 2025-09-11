// Package conf return the embed files for configuration.
package conf

import "embed"

//go:embed *.toml
var FileFS embed.FS
