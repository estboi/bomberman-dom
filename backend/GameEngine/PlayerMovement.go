package gameengine

import (
	c "bomberman/Constants"
	"log"
	"math"
)

type PlayerCollision struct {
	Axis           string
	Pickup         *Powerup
	DistanceToEdge float64
	EdgeAxis       int
	CurrentPos     c.Coords
	TileToMove     *c.Coords

	player *Player
}

func (Game *GameState) PlayerInput(input c.PlayerInput) {
	player, ok := Game.Players[input.PlayerID]
	if !ok {
		return
	}

	player.Direction = input.Direction
	var X = player.Position.X
	var Y = player.Position.Y

	var positionToMove = PlayerCollision{}

	switch input.Direction {
	case "U":
		positionToMove = PlayerCollision{
			Axis:       "-Y",
			TileToMove: &c.Coords{X: X + c.HITBOX_CENTER, Y: Y - player.Stats.Speed},
		}
	case "D":
		positionToMove = PlayerCollision{
			Axis:       "Y",
			TileToMove: &c.Coords{X: X + c.HITBOX_CENTER, Y: Y + player.Stats.Speed + c.PLAYER_HITBOX_SIZE},
		}
	case "R":
		positionToMove = PlayerCollision{
			Axis:       "X",
			TileToMove: &c.Coords{X: X + player.Stats.Speed + c.PLAYER_HITBOX_SIZE, Y: Y + c.HITBOX_CENTER},
		}
	case "L":
		positionToMove = PlayerCollision{
			Axis:       "-X",
			TileToMove: &c.Coords{X: X - player.Stats.Speed, Y: Y + c.HITBOX_CENTER},
		}
	case "P":
		Game.PlacingBomb(player)
		return
	}
	positionToMove.CurrentPos = player.Position
	positionToMove.player = player

	if Game.collisionMechanics(&positionToMove) {
		player.Hit(Game.Updates)
	}
	player.Move(positionToMove, Game)
}

func (Game *GameState) collisionMechanics(positionToMove *PlayerCollision) bool {
	var TileToMove = c.Coords{X: math.Floor(c.Round(positionToMove.TileToMove.X)), Y: math.Floor(c.Round(positionToMove.TileToMove.Y))}
	obj, isCollision := Game.ColladibleObjects[TileToMove]
	if obj.ID == 0 && obj.Type == 0 {
		isCollision = false
	}
	if !isCollision {
		// Check if player hitbox fits
		switch positionToMove.Axis {
		case "X", "-X":
			if positionToMove.CurrentPos.Y+c.PLAYER_HITBOX_SIZE > TileToMove.Y+1 {
				positionToMove.DistanceToEdge = c.Round(positionToMove.CurrentPos.Y + c.PLAYER_HITBOX_SIZE - (TileToMove.Y + 1))
				positionToMove.EdgeAxis = -1
			} else if positionToMove.CurrentPos.Y > TileToMove.Y-1 && positionToMove.CurrentPos.Y < TileToMove.Y {
				positionToMove.DistanceToEdge = c.Round(TileToMove.Y - positionToMove.CurrentPos.Y)
				positionToMove.EdgeAxis = 1
			}
		case "Y", "-Y":
			if positionToMove.CurrentPos.X+c.PLAYER_HITBOX_SIZE > TileToMove.X+1 {
				positionToMove.DistanceToEdge = c.Round(positionToMove.CurrentPos.X + c.PLAYER_HITBOX_SIZE - (TileToMove.X + 1))
				positionToMove.EdgeAxis = -1
			} else if positionToMove.CurrentPos.X > TileToMove.X-1 && positionToMove.CurrentPos.X < TileToMove.X {
				positionToMove.DistanceToEdge = c.Round(TileToMove.X - positionToMove.CurrentPos.X)
				positionToMove.EdgeAxis = 1
			}
		}
		positionToMove.TileToMove = nil
		return false
	}

	switch obj.Type {
	case c.WALL, c.CRATE, c.BOMB:
		positionToMove.TileToMove = &c.Coords{X: math.Floor(positionToMove.TileToMove.X), Y: math.Floor(positionToMove.TileToMove.Y)}
		return false
	case c.POWERUP_SPEED, c.POWERUP_BOMB, c.POWERUP_HEAL, c.POWERUP_POWER:
		Game.PowerupCollision(obj, positionToMove)
		Game.ColladibleObjects[TileToMove] = ColladibleObject{}
		positionToMove.TileToMove = nil
		return false
	case c.BLAST:
		if !positionToMove.player.IsIFrame {
			return true
		}
	}
	positionToMove.TileToMove = nil
	return false
}

func (Game *GameState) PowerupCollision(powerup ColladibleObject, collision *PlayerCollision) {
	powerupToPickUp := Game.Powerups[powerup.ID]
	collision.Pickup = powerupToPickUp
	powerupToPickUp.Death(Game.Updates)
	powerupToPickUp = nil

}

func (Game *GameState) PlacingBomb(player *Player) {
	if player.PlacedBombs <= player.Stats.Bombs {
		if Game.isAnotherBombOnTile(player) {
			log.Println("There is another bomb on this tile")
			return
		}
		var newBomb Bomb
		newBomb.Initialize(player.Position, player.Stats.Power, player.ID, Game.Updates)
		Game.Bombs[newBomb.ID] = &newBomb
		Game.ColladibleObjects[newBomb.Position] = ColladibleObject{ID: newBomb.ID, Type: c.BOMB}
		player.PlacedBombs++
	}
}

func (Game *GameState) isAnotherBombOnTile(player *Player) bool {
	X := math.Floor(player.Position.X + 0.5)
	Y := math.Floor(player.Position.Y + 0.5)
	obj, isExists := Game.ColladibleObjects[c.Coords{X: X, Y: Y}]
	if obj.ID == 0 && obj.Type == 0 {
		return false
	}
	log.Println(obj)
	return isExists
}
