//go:build ignore

package main

import (
	"log/slog"
	"time"

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

	mdk.GetModEventBus().Register("EverySecond", func(arg interface{}) {
		slog.Info("[mod] EverySecond triggered", "time", arg.(time.Time))
	})
	return nil
}

func Close() error {
	return nil
}

func Modify(game *mdk.Game) error {
	game.Horalky = 42

	return nil
}
