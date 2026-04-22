package main

import (
	"github.com/basel-ax/flying-pngs/internal/config"
	"github.com/basel-ax/flying-pngs/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cfg, _ := config.Load()
	game := ui.NewGame(&cfg)
	ebiten.RunGame(game)
}
