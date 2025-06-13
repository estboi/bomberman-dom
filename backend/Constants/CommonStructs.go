package constants

import "math"

type Coords struct {
	X float64
	Y float64
}

var PlayerSpawn = map[int]Coords{
	11: {X: 1, Y: 1},
	12: {X: 14, Y: 1},
	13: {X: 1, Y: 14},
	14: {X: 14, Y: 14},
}

type PlayerInput struct {
	Direction string
	PlayerID  int
}

type BlastCoords struct {
	Center Coords
	Up     []Coords
	Right  []Coords
	Down   []Coords
	Left   []Coords
}

func Round(number float64) float64 {
	return math.Round(number*1e6) / 1e6
}
