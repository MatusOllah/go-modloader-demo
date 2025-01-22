//go:build ignore

package main

import (
	"image/color"
	"log/slog"

	"github.com/MatusOllah/go-modloader-demo/mdk"
)

func Metadata() *mdk.ModMetadata {
	return &mdk.ModMetadata{
		Name:        "Hello Mod",
		Version:     "1.0.0",
		Author:      "MatusOllah",
		Description: "A test mod.",
	}
}

func Init() error {
	slog.Info("[mod] Hello from Init func!")

	bus := mdk.GetModEventBus()

	bus.Register("OnDeath", func(args interface{}) {
		slog.Info("[mod] OnDeath triggered", "args", args.(*mdk.OnDeathEventArgs))
	})
	bus.Register("OnAppleCollected", func(args interface{}) {
		slog.Info("[mod] OnAppleCollected triggered", "args", args.(*mdk.OnAppleCollectedEventArgs))
	})
	return nil
}

func Close() error {
	return nil
}

func Modify(game *mdk.Game) error {
	slog.Info("[mod] Hello from Modify func!")
	game.SnakeColor = color.RGBA{0, 255, 0, 255}

	return nil
}
