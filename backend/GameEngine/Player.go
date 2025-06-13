package gameengine

import (
	c "bomberman/Constants"
	"encoding/json"
	"log"
	"math"
	"time"
)

type Player struct {
	ID       int
	Position c.Coords
	Stats    PlayerStats

	PlacedBombs int
	IsIFrame    bool
	IFrameTimer time.Time

	IsAlive   bool
	Direction string
}

type PlayerStats struct {
	Health int
	Speed  float64
	Power  int
	Bombs  int
}

type PlayerMoveUpdate struct {
	PlayerID  int
	Coords    c.Coords
	Direction string
}

type PlayerHitUpdate struct {
	PlayerID int
	Health   int
}

type PlayerPowerupUpdate struct {
	ID   int
	Type int
}

func (p *Player) Initialize(playerID int) {
	p.ID = playerID
	p.Stats = PlayerStats{
		Health: 3,
		Speed:  0.1,
		Power:  1,
		Bombs:  1,
	}
	p.IFrameTimer = time.Time{}
	p.IsAlive = true
}

func (p *Player) Move(moveInput PlayerCollision, Game *GameState) {
	distanceToMove := p.Stats.Speed

	log.Println(p.Stats.Health)

	if moveInput.Pickup != nil {
		p.Pickup(moveInput.Pickup, Game.Updates)
	}
	if moveInput.TileToMove != nil {
		// Caluclate the distance till collision
		var deltaDistance float64

		switch moveInput.Axis {
		case "X":
			deltaDistance = math.Mod(moveInput.TileToMove.X-(p.Position.X+c.PLAYER_HITBOX_SIZE), 1)
		case "-X":
			deltaDistance = math.Mod(moveInput.TileToMove.X-p.Position.X, 1)
		case "Y":
			deltaDistance = math.Mod(moveInput.TileToMove.Y-(p.Position.Y+c.PLAYER_HITBOX_SIZE), 1)
		case "-Y":
			deltaDistance = math.Mod(moveInput.TileToMove.Y-p.Position.Y, 1)
		}

		distanceToMove = math.Round(math.Abs(deltaDistance)*10) / 10

		switch {
		case distanceToMove > p.Stats.Speed:
			distanceToMove -= (distanceToMove - p.Stats.Speed)
		case distanceToMove == 0:
			return
		}
	}

	var newPos float64

	switch moveInput.Axis {
	case "X":
		newPos = c.Round(distanceToMove + p.Position.X)
		p.Position.Y += p.calculateDistanceToEdge(moveInput.EdgeAxis, moveInput.DistanceToEdge)
		p.Position.X = math.Abs(newPos)
	case "-X":
		newPos = c.Round(distanceToMove - p.Position.X)
		p.Position.Y += p.calculateDistanceToEdge(moveInput.EdgeAxis, moveInput.DistanceToEdge)
		p.Position.X = math.Abs(newPos)
	case "Y":
		newPos = c.Round(distanceToMove + p.Position.Y)
		p.Position.X += p.calculateDistanceToEdge(moveInput.EdgeAxis, moveInput.DistanceToEdge)
		p.Position.Y = math.Abs(newPos)
	case "-Y":
		newPos = c.Round(distanceToMove - p.Position.Y)
		p.Position.X += p.calculateDistanceToEdge(moveInput.EdgeAxis, moveInput.DistanceToEdge)
		p.Position.Y = math.Abs(newPos)
	}

	Game.Updates <- p.Update("move", nil)
}

func (p *Player) calculateDistanceToEdge(axis int, distanceToEdge float64) float64 {
	if p.Stats.Speed == 0.1 {
		if distanceToEdge-0.1 >= 0 {
			return float64(axis) * 0.1
		}
	}
	return float64(axis) * distanceToEdge
}

func (p *Player) Pickup(pickupType *Powerup, updateChan chan UpdatedEntity) {
	switch pickupType.Type {
	case c.POWERUP_SPEED:
		if p.Stats.Speed <= 0.3 {
			p.Stats.Speed += c.SPEED_POWERUP_MULTIPLIER
		}
	case c.POWERUP_BOMB:
		if p.Stats.Bombs <= 7 {
			p.Stats.Bombs += 1
		}
	case c.POWERUP_POWER:
		if p.Stats.Power <= 5 {
			p.Stats.Power += 1
		}
	case c.POWERUP_HEAL:
		if p.Stats.Health != 3 {
			p.Stats.Health += 1
		}
	}
	updateChan <- p.Update("player_powerup", pickupType)
}

func (p *Player) Hit(update chan UpdatedEntity) {
	p.Stats.Health--
	if p.Stats.Health == 0 {
		p.Death(update)
	} else {
		p.IsIFrame = true
		p.IFrameTimer = time.Now()
		update <- p.Update("player_hit", nil)
	}
}

func (p *Player) Death(update chan UpdatedEntity) {
	update <- p.Update("player_dead", nil)
}

func (p *Player) Update(eventType string, powerup *Powerup) UpdatedEntity {
	var payload interface{}

	switch eventType {
	case "move":
		payload = PlayerMoveUpdate{
			PlayerID:  p.ID,
			Coords:    p.Position,
			Direction: p.Direction,
		}
	case "player_hit":
		payload = PlayerHitUpdate{
			PlayerID: p.ID,
			Health:   p.Stats.Health,
		}
	case "player_dead":
		payload = p.ID
	case "player_powerup":
		payload = PlayerPowerupUpdate{
			ID:   powerup.ID,
			Type: powerup.Type,
		}
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
