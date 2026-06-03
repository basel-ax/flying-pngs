package main

import (
	"log"
	"os"

	"github.com/basel-ax/flying-pngs/internal/config"
	"github.com/basel-ax/flying-pngs/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("[WARN] Could not load config, using defaults: %v", err)
		cfg = config.DefaultConfig()
	}

	game := ui.NewGame(&cfg)

	ebiten.SetWindowSize(cfg.Width, cfg.Height)
	ebiten.SetWindowTitle("Flying PNGs — Settings")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	opts := &ebiten.RunGameOptions{
		ScreenTransparent: cfg.BackgroundType == "transparent",
	}
	if err := ebiten.RunGameWithOptions(game, opts); err != nil {
		log.Printf("[ERROR] Game exited: %v", err)
		os.Exit(1)
	}
}
