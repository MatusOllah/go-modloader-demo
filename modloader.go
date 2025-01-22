package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/MatusOllah/go-modloader-demo/mdk"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

var Symbols interp.Exports = make(interp.Exports)

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
	i.Use(Symbols)
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
		return fmt.Errorf("failed to read mod bytes: %w", err)
	}

	_, err = i.Eval(string(src))
	if err != nil {
		return fmt.Errorf("failed to compile mod: %w", err)
	}

	metadata := getSym(i, "main.Metadata").Interface().(func() *mdk.ModMetadata)()
	slog.Info("got metadata", "metadata", metadata)

	if err := getSym(i, "main.Init").Interface().(func() error)(); err != nil {
		return fmt.Errorf("failed to initialize mod: %w", err)
	}

	mdkGame := &mdk.Game{
		SnakeColor: g.snakeColor,
		AppleColor: g.appleColor,
	}

	if err := getSym(i, "main.Modify").Interface().(func(*mdk.Game) error)(mdkGame); err != nil {
		return fmt.Errorf("failed to modify game: %w", err)
	}

	g.snakeColor = mdkGame.SnakeColor
	g.appleColor = mdkGame.AppleColor

	return nil
}

func getSym(i *interp.Interpreter, name string) reflect.Value {
	v, err := i.Eval(name)
	if err != nil {
		panic(err)
	}

	return v
}
