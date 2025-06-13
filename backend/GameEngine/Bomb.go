package gameengine

import (
	c "bomberman/Constants"
	"encoding/json"
	"log"
	"math"
	"time"
)

type Bomb struct {
	ID           int
	Power        int
	Position     c.Coords
	PlayerPlaced int

	StartTimer time.Time
}

type BombUpdate struct {
	ID     int
	Coords c.Coords
}

type Blast struct {
	ID          int
	Position    c.Coords
	BlastCoords c.BlastCoords
	Power       int

	StartTimer time.Time
}

type BlastUpdate struct {
	ID     int
	Coords c.BlastCoords
}

func (b *Bomb) Initialize(spawnPoint c.Coords, power int, playerId int, update chan UpdatedEntity) {
	X := math.Floor(spawnPoint.X + 0.5)
	Y := math.Floor(spawnPoint.Y + 0.5)

	b.Position = c.Coords{X: X, Y: Y}
	b.ID = createID(b.Position)
	b.Power = power
	b.PlayerPlaced = playerId

	b.StartTimer = time.Now()
	update <- b.Update("bomb")
}

func createID(pos c.Coords) int {
	var id int
	id += int(pos.Y)*100 + int(pos.X)
	return id
}

func (b *Bomb) Update(eventType string) UpdatedEntity {
	payload := map[string]interface{}{
		"ID": b.ID,
		"Coords": map[string]float64{
			"X": b.Position.X,
			"Y": b.Position.Y,
		},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		// Handle JSON marshaling error
		log.Println("Error marshaling JSON:", err)
		// Returning an empty UpdatedEntity as an indication of an error
		return UpdatedEntity{}
	}

	update := UpdatedEntity{
		Type:    eventType,
		Payload: payloadJSON,
	}
	return update
}

func (b *Blast) Initialize(bomb Bomb) {
	b.ID = bomb.ID
	b.Position = bomb.Position
	b.Power = bomb.Power
}

func (b *Blast) Death(update chan UpdatedEntity) {
	update <- b.Update("blast_dead")
}

func (b *Blast) Update(eventType string) UpdatedEntity {
	var payload = BlastUpdate{
		ID:     b.ID,
		Coords: b.BlastCoords,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling JSON: ", err)
		return UpdatedEntity{}
	}
	update := UpdatedEntity{
		Type:    eventType,
		Payload: payloadJSON,
	}
	return update
}
