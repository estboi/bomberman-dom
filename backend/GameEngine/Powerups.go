package gameengine

import (
	c "bomberman/Constants"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
)

type Powerup struct {
	ID       int
	Type     int
	Position c.Coords
	IsAlive  bool
}

type PowerupUpdate struct {
	ID    int
	Type  int
	Coord c.Coords
}

func (P *Powerup) Initialize(Position c.Coords, update chan UpdatedEntity) error {

	P.IsAlive = true

	// Calculate powerup type based on spawn chance
	// If didn't met the spanw chance don't add to gameloop
	P.GenerateType()
	if !P.IsAlive {
		return errors.New("didn't met the spawn chance")
	}
	P.Position = Position
	P.CreateID()

	update <- P.Update("powerup")

	return nil
}

func (P *Powerup) CreateID() {
	posSum := P.Position.Y*100 + P.Position.X
	P.ID = int(posSum)
}

func (P *Powerup) GenerateType() {
	randomInt := rand.Intn(100)
	if randomInt <= c.POWERUP_DROP_CHANCE {
		switch {
		case randomInt <= c.HEAL_DROP_CHANCE:
			P.Type = c.POWERUP_HEAL
		case randomInt > c.HEAL_DROP_CHANCE && randomInt <= c.BOMB_DROP_CHANCE+c.HEAL_DROP_CHANCE:
			P.Type = c.POWERUP_BOMB
		case randomInt > c.BOMB_DROP_CHANCE+c.HEAL_DROP_CHANCE && randomInt <= c.POWER_DROP_CHANCE+c.BOMB_DROP_CHANCE:
			P.Type = c.POWERUP_POWER
		case randomInt > c.SPEED_DROP_CHANCE+c.HEAL_DROP_CHANCE:
			P.Type = c.POWERUP_SPEED
		default:
			P.IsAlive = false
		}
		return
	}
	P.IsAlive = false
}

func (p *Powerup) Death(update chan UpdatedEntity) {
	p.IsAlive = false
	log.Println(p)
	update <- p.Update("powerup_dead")
}

func (p *Powerup) Update(eventType string) UpdatedEntity {
	var payload interface{}

	switch eventType {
	case "powerup":
		payload = PowerupUpdate{
			ID:    p.ID,
			Type:  p.Type,
			Coord: p.Position,
		}
	case "powerup_dead":
		payload = p.ID
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return UpdatedEntity{}
	}

	update := UpdatedEntity{
		Type:    eventType,
		Payload: payloadJSON,
	}
	return update
}
