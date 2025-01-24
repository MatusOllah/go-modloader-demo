//go:build ignore

package main

import (
	"log/slog"

	"github.com/MatusOllah/go-modloader-demo/mdk"
)

func Metadata() *mdk.ModMetadata {
	return &mdk.ModMetadata{
		ID:          "hello",
		DisplayName: "Hello Mod",
		Version:     "1.0.0",
		Author:      "MatusOllah",
		License:     "MIT License",
		Description: "A test mod.",
	}
}

func Init() error {
	slog.Info("[mod] Hello from Init func!")

	bus := mdk.ModEventBus

	bus.Register("OnDeath", func(args interface{}) {
		slog.Info("[mod] OnDeath triggered", "args", args.(*mdk.OnDeathEventArgs))
	})
	bus.Register("OnAppleCollected", func(args interface{}) {
		slog.Info("[mod] OnAppleCollected triggered", "args", args.(*mdk.OnAppleCollectedEventArgs))
	})

	mdk.ThingRegistry.Register("lamp", mdk.Thing("Lamp"))
	mdk.ThingRegistry.Register("microphone", mdk.Thing("Microphone"))
	return nil
}

func Close() error {
	slog.Info("[mod] Goodbye!")

	return nil
}
