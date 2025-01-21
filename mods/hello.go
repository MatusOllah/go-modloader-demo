//go:build ignore

package main

import "github.com/MatusOllah/go-modloader-demo/mdk"

func Metadata() *mdk.ModMetadata {
	return &mdk.ModMetadata{
		Name:        "Hello Mod",
		Version:     "1.0.0",
		Author:      "MatusOllah",
		Description: "A test mod.",
	}
}

func Init() error {
	return nil
}

func Close() error {
	return nil
}

func Modify(game *mdk.Game) error {
	game.Horalky = 42

	return nil
}
