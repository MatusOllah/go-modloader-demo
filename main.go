package main

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/MatusOllah/go-modloader-demo/mdk"
	"github.com/MatusOllah/slogcolor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/traefik/yaegi/interp"
)

const (
	screenWidth        int = 640
	screenHeight       int = 480
	gridSize               = 10
	xGridCountInScreen int = screenWidth / gridSize
	yGridCountInScreen int = screenHeight / gridSize
)

type Direction int

const (
	DirectionNone Direction = iota
	DirectionUp
	DirectionDown
	DirectionLeft
	DirectionRight
)

type Game struct {
	moveDirection Direction
	snakeBody     []image.Point
	apple         image.Point
	timer         int
	moveTime      int
	score         int
	bestScore     int
	level         int

	snakeColor color.Color
	appleColor color.Color
}

func NewGame() *Game {
	g := &Game{
		apple:      image.Pt(3*gridSize, 3*gridSize),
		moveTime:   4,
		snakeBody:  make([]image.Point, 1),
		snakeColor: color.RGBA{0x80, 0xa0, 0xc0, 0xff},
		appleColor: color.RGBA{0xFF, 0x00, 0x00, 0xff},
	}
	g.snakeBody[0].X = xGridCountInScreen / 2
	g.snakeBody[0].Y = yGridCountInScreen / 2
	return g
}

func (g *Game) collidesWithApple() bool {
	return g.snakeBody[0].X == g.apple.X && g.snakeBody[0].Y == g.apple.Y
}

func (g *Game) collidesWithSelf() bool {
	for _, v := range g.snakeBody[1:] {
		if g.snakeBody[0].X == v.X &&
			g.snakeBody[0].Y == v.Y {
			return true
		}
	}
	return false
}

func (g *Game) collidesWithWall() bool {
	return g.snakeBody[0].X < 0 ||
		g.snakeBody[0].Y < 0 ||
		g.snakeBody[0].X >= xGridCountInScreen ||
		g.snakeBody[0].Y >= yGridCountInScreen
}

func (g *Game) needsToMoveSnake() bool {
	return g.timer%g.moveTime == 0
}

func (g *Game) reset() {
	g.apple.X = 3 * gridSize
	g.apple.Y = 3 * gridSize
	g.moveTime = 4
	g.snakeBody = g.snakeBody[:1]
	g.snakeBody[0].X = xGridCountInScreen / 2
	g.snakeBody[0].Y = yGridCountInScreen / 2
	g.score = 0
	g.level = 1
	g.moveDirection = DirectionNone
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if g.moveDirection != DirectionRight {
			g.moveDirection = DirectionLeft
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if g.moveDirection != DirectionLeft {
			g.moveDirection = DirectionRight
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.moveDirection != DirectionUp {
			g.moveDirection = DirectionDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.moveDirection != DirectionDown {
			g.moveDirection = DirectionUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.reset()
	}

	if g.needsToMoveSnake() {
		if g.collidesWithWall() || g.collidesWithSelf() {
			g.reset()
			mdk.ModEventBus.Trigger("OnDeath", &mdk.OnDeathEventArgs{
				SnakeBody: g.snakeBody,
				Level:     g.level,
				Score:     g.score,
				BestScore: g.bestScore,
			})
		}

		if g.collidesWithApple() {
			oldApple := g.apple

			g.apple.X = rand.IntN(xGridCountInScreen - 1)
			g.apple.Y = rand.IntN(yGridCountInScreen - 1)
			g.snakeBody = append(g.snakeBody, image.Point{
				X: g.snakeBody[len(g.snakeBody)-1].X,
				Y: g.snakeBody[len(g.snakeBody)-1].Y,
			})
			if len(g.snakeBody) > 10 && len(g.snakeBody) < 20 {
				g.level = 2
				g.moveTime = 3
			} else if len(g.snakeBody) > 20 {
				g.level = 3
				g.moveTime = 2
			} else {
				g.level = 1
			}
			g.score++
			if g.bestScore < g.score {
				g.bestScore = g.score
			}
			mdk.ModEventBus.Trigger("OnAppleCollected", &mdk.OnAppleCollectedEventArgs{
				SnakeBody: g.snakeBody,
				OldApple:  oldApple,
				NewApple:  g.apple,
				Level:     g.level,
				Score:     g.score,
				BestScore: g.bestScore,
			})
		}

		for i := int64(len(g.snakeBody)) - 1; i > 0; i-- {
			g.snakeBody[i].X = g.snakeBody[i-1].X
			g.snakeBody[i].Y = g.snakeBody[i-1].Y
		}
		switch g.moveDirection {
		case DirectionLeft:
			g.snakeBody[0].X--
		case DirectionRight:
			g.snakeBody[0].X++
		case DirectionDown:
			g.snakeBody[0].Y++
		case DirectionUp:
			g.snakeBody[0].Y--
		}
	}

	g.timer++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, v := range g.snakeBody {
		vector.DrawFilledRect(screen, float32(v.X*gridSize), float32(v.Y*gridSize), gridSize, gridSize, g.snakeColor, false)
	}
	vector.DrawFilledRect(screen, float32(g.apple.X*gridSize), float32(g.apple.Y*gridSize), gridSize, gridSize, g.appleColor, false)

	if g.moveDirection == DirectionNone {
		ebitenutil.DebugPrintAt(screen, "Press up/down/left/right to start", 200, 200)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTPS: %.2f\nLevel: %d\nScore: %d\nBest Score: %d\n", ebiten.ActualFPS(), ebiten.ActualTPS(), g.level, g.score, g.bestScore))
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, nil)))

	slog.Info("initializing game")
	g := NewGame()

	slog.Info("loading mods")
	mods, err := filepath.Glob("mods/*.go")
	if err != nil {
		slog.Error("failed to scan mods", "err", err)
	}

	if err := loadMods(g, &ModloaderOptions{Mods: mods, InterpOpts: &interp.Options{Unrestricted: true}}); err != nil {
		slog.Error("failed to load mods", "err", err)
		os.Exit(1)
	}
	defer closeMods()

	slog.Info("initalizing ebiten")
	ebiten.SetWindowTitle("Go Yaegi Modloader Demo")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	slog.Info("running game")
	if err := ebiten.RunGame(g); err != nil {
		slog.Error("failed to run game", "err", err)
		os.Exit(1)
	}

}
