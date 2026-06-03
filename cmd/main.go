package main

import (
	"flag"
	"os"

	"github.com/basel-ax/flying-pngs/internal/config"
	"github.com/basel-ax/flying-pngs/internal/logger"
	"github.com/basel-ax/flying-pngs/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode with verbose [INFO] logging")
	flag.Parse()

	logger.SetDebug(*debug)
	if *debug {
		logger.Info("Debug mode enabled")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Warn("Could not load config, using defaults: %v", err)
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
		logger.Error("Game exited: %v", err)
		os.Exit(1)
	}
}
