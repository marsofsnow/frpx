package frpc

import (
	"embed"

	"github.com/marsofsnow/frpx/assets"
)

//go:embed static/*
var content embed.FS

func init() {
	assets.Register(content)
}
