package main

import (
	"log/slog"
	"os"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

var Symbols interp.Exports

type ModloaderOptions struct {
	Mods         []string
	InterpOpts   *interp.Options
	UseSyscall   bool
	UseUnsafe    bool
	ExtraSymbols interp.Exports
}

// TODO: this
func loadMods(g *Game, opts *ModloaderOptions) error {
	for _, path := range opts.Mods {
		slog.Info("loading mod", "path", path)
		if err := loadMod(g, path, opts); err != nil {
			return err
		}
	}

	return nil
}

func loadMod(g *Game, path string, opts *ModloaderOptions) error {
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

	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = i.Eval(string(src))
	if err != nil {
		return err
	}

	return nil
}

func getSym(i *interp.Interpreter, name string) reflect.Value {
	v, err := i.Eval(name)
	if err != nil {
		panic(err)
	}

	return v
}
