package resources

import (
	"embed"
)

//go:embed uaa.yml
var UaaConfigFs embed.FS

var UaaConfig, _ = UaaConfigFs.ReadFile("uaa.yml")
