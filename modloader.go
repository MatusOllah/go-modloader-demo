package main

import (
	"os"

	"github.com/traefik/yaegi/interp"
)

// TODO: this
func loadMods(g *Game) error {
	i := interp.New(interp.Options{
		Env:          os.Environ(),
		Unrestricted: true,
	})

	return nil
}
