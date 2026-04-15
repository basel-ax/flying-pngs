package main

import (
	"flying-pngs-go/internal/config"
	"flying-pngs-go/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cfg, _ := config.Load()
	game := ui.NewGame(&cfg)
	ebiten.RunGame(game)
}
