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

type Mod struct {
	Metadata *mdk.ModMetadata
	i        *interp.Interpreter
}

func (m *Mod) Close() error {
	v, err := m.i.Eval("main.Close")
	if err != nil {
		return err
	}
	return v.Interface().(func() error)()
}

var Mods map[string]*Mod = make(map[string]*Mod)

type ModloaderOptions struct {
	Mods         []string
	InterpOpts   *interp.Options
	UseSyscall   bool
	UseUnsafe    bool
	ExtraSymbols interp.Exports
}

// TODO: make Game struct moddable
func loadMods(g *Game, opts *ModloaderOptions) error {
	for _, path := range opts.Mods {
		slog.Info("loading mod", "path", path)
		m, err := loadMod(path, opts)
		if err != nil {
			return err
		}
		Mods[m.Metadata.ID] = m
	}

	return nil
}

func loadMod(path string, opts *ModloaderOptions) (*Mod, error) {
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
		return nil, fmt.Errorf("failed to read mod bytes: %w", err)
	}

	_, err = i.Eval(string(src))
	if err != nil {
		return nil, fmt.Errorf("failed to compile mod: %w", err)
	}

	metadata := getSym(i, "main.Metadata").Interface().(func() *mdk.ModMetadata)()
	slog.Info("got metadata", "metadata", metadata)

	if err := getSym(i, "main.Init").Interface().(func() error)(); err != nil {
		return nil, fmt.Errorf("failed to initialize mod: %w", err)
	}

	return &Mod{
		Metadata: metadata,
		i:        i,
	}, nil
}

func getSym(i *interp.Interpreter, name string) reflect.Value {
	v, err := i.Eval(name)
	if err != nil {
		panic(err)
	}

	return v
}

func closeMods() {
	for id, mod := range Mods {
		slog.Info("closing mod", "id", id)
		if err := mod.Close(); err != nil {
			panic(err)
		}
	}
}
