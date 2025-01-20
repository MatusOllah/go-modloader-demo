package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/MatusOllah/slogcolor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	screenWidth  int = 640
	screenHeight int = 480
)

type Game struct {
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f; TPS: %.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, nil)))

	slog.Info("initializing game")
	g := NewGame()

	slog.Info("initalizing ebiten")
	ebiten.SetWindowTitle("Go Yaegi Modloader Demo")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	slog.Info("running game")
	if err := ebiten.RunGame(g); err != nil {
		slog.Error("failed to run game", "err", err)
		os.Exit(1)
	}

}
