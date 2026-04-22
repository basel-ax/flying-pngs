package ui

import (
	"github.com/basel-ax/flying-pngs/internal/animation"
	"github.com/basel-ax/flying-pngs/internal/config"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	cfg  *config.Config
	anim *animation.AnimationCanvas
}

func NewGame(cfg *config.Config) *Game {
	g := &Game{
		cfg:  cfg,
		anim: animation.NewAnimationCanvas(cfg),
	}
	g.anim.Start()
	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.Width, g.cfg.Height
}

func (g *Game) Update() error {
	g.anim.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.anim.Draw(screen)
}
