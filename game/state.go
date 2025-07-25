package game

type GameState int

const (
	MenuState GameState = iota
	PlayingState
	PauseState
)