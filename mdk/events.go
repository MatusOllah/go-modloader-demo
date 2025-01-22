package mdk

import "image"

type OnDeathEventArgs struct {
	SnakeBody []image.Point
	Level     int
	Score     int
	BestScore int
}

type OnAppleCollectedEventArgs struct {
	SnakeBody []image.Point
	OldApple  image.Point
	NewApple  image.Point
	Level     int
	Score     int
	BestScore int
}
