package main

import (
	"embed"

	"github.com/lostsnow/cloudrain/cmd"
)

//go:embed web/dist
var webAssets embed.FS

func main() {
	cmd.Execute(&cmd.Assets{Web: webAssets})
}
