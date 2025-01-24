package mdk

import "image/color"

// GameInstance is the current game instance.
var GameInstance *Game

type Game struct {
	SnakeColor color.Color
	AppleColor color.Color
}
