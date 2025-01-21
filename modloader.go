package main

import (
	"os"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

type ModloaderOptions struct {
	Mods         []string
	InterpOpts   *interp.Options
	UseSyscall   bool
	UseUnsafe    bool
	ExtraSymbols interp.Exports
}

// TODO: this
func loadMods(g *Game, opts *ModloaderOptions) error {
	i := interp.New(*opts.InterpOpts)

	i.Use(stdlib.Symbols)
	i.Use(interp.Symbols)
	if opts.UseSyscall {
		i.Use(syscall.Symbols)
	}
	if opts.UseUnsafe {
		i.Use(unsafe.Symbols)
	}
	if opts.ExtraSymbols != nil {
		i.Use(opts.ExtraSymbols)
	}

	for _, path := range opts.Mods {
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		p, err := i.Compile(string(src))
		if err != nil {
			return err
		}

		_, err = i.Execute(p)
		if err != nil {
			return err
		}
	}

	return nil
}
