package gameengine

import (
	c "bomberman/Constants"
	"encoding/json"
	"sync"
)

type GameState struct {
	Map MapData

	Crates   map[int]*Crate
	Players  map[int]*Player
	Bombs    map[int]*Bomb
	Blast    map[int]*Blast
	Powerups map[int]*Powerup

	ColladibleObjects map[c.Coords]ColladibleObject
	Mutex             sync.Mutex

	Updates chan UpdatedEntity
}

type ColladibleObject struct {
	ID   int
	Type int
}

type UpdatedEntity struct {
	Type    string
	Payload json.RawMessage
}

func (Game *GameState) Initialize() {
	Game.ColladibleObjects = make(map[c.Coords]ColladibleObject)
	Game.Crates = make(map[int]*Crate)
	Game.Players = make(map[int]*Player)
	Game.Bombs = make(map[int]*Bomb)
	Game.Blast = make(map[int]*Blast)
	Game.Powerups = make(map[int]*Powerup)

	Game.Updates = make(chan UpdatedEntity, 100)

	Game.InitializeCollidible()
}

func (Game *GameState) InitializeCollidible() {
	Game.Map.Initialize()
	Game.AddCrates()
	Game.AddWalls()
}

func (Game *GameState) AddCrates() {
	Game.Crates = Game.Map.Crates
	for _, crate := range Game.Crates {
		Game.ColladibleObjects[crate.Position] = ColladibleObject{ID: crate.ID, Type: crate.Type}
	}
}

func (Game *GameState) AddWalls() {
	matrix := Game.Map.MapMatrix
	for y, row := range matrix {
		for x, obj := range row {
			if obj != c.WALL {
				continue
			}
			Game.ColladibleObjects[c.Coords{X: float64(x), Y: float64(y)}] = ColladibleObject{ID: obj, Type: obj}
		}
	}
}

// Call this function when new user has logged in
func (Game *GameState) InitializePlayer(playerID int) {
	Coords := Game.Map.MapMatrix.SpawnPlayer(playerID)
	Game.Players[playerID] = &Player{Position: Coords}
	Game.Players[playerID].Initialize(playerID)
}

func (Game *GameState) ClearBlast(blast *Blast) {
	delete(Game.Blast, blast.ID)
	for _, coords := range blast.BlastCoords.Left {
		Game.ColladibleObjects[coords] = ColladibleObject{}
	}
	for _, coords := range blast.BlastCoords.Up {
		Game.ColladibleObjects[coords] = ColladibleObject{}
	}
	for _, coords := range blast.BlastCoords.Down {
		Game.ColladibleObjects[coords] = ColladibleObject{}
	}
	for _, coords := range blast.BlastCoords.Right {
		Game.ColladibleObjects[coords] = ColladibleObject{}
	}
	Game.ColladibleObjects[blast.BlastCoords.Center] = ColladibleObject{}
}

func (Game *GameState) ClearPowerup(powerup PowerupUpdate) {
	Game.Powerups[powerup.ID] = nil
	Game.ColladibleObjects[powerup.Coord] = ColladibleObject{}
}
