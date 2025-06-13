package gameengine

import (
	c "bomberman/Constants"
	"fmt"
	"math"
	"time"
)

func (Game *GameState) Explode(bomb *Bomb) {
	var newBlast Blast
	newBlast.Initialize(*bomb)

	Game.Players[bomb.PlayerPlaced].PlacedBombs -= 1

	Game.ColladibleObjects[bomb.Position] = ColladibleObject{}
	Game.Blast[newBlast.ID] = &newBlast
	delete(Game.Bombs, bomb.ID)

	Game.CalculateBlast(&newBlast)
}

func (Game *GameState) CalculateBlast(blast *Blast) {
	var blastCoords = c.BlastCoords{Center: blast.Position}
	Game.ColladibleObjects[blastCoords.Center] = ColladibleObject{ID: blast.ID, Type: c.BLAST}
	for _, player := range Game.Players {
		if !player.IsIFrame {
			if math.Floor(player.Position.Y) == blast.Position.X && math.Floor(player.Position.Y) == blast.Position.Y {
				player.Hit(Game.Updates)
			}
		}
	}

	blastCoords.Up = Game.checkBlastCollision(blast, 0, -1)
	blastCoords.Right = Game.checkBlastCollision(blast, 1, 0)
	blastCoords.Down = Game.checkBlastCollision(blast, 0, 1)
	blastCoords.Left = Game.checkBlastCollision(blast, -1, 0)
	blast.BlastCoords = blastCoords
	blast.StartTimer = time.Now()

	Game.Updates <- blast.Update("blast")
}

func (Game *GameState) checkBlastCollision(blast *Blast, deltaX, deltaY int) []c.Coords {
	var blastDirection = []c.Coords{}
	for power := 1; power <= blast.Power; power++ {
		nextX := blast.Position.X + float64(deltaX*power)
		nextY := blast.Position.Y + float64(deltaY*power)
		var nextCoords = c.Coords{X: nextX, Y: nextY}
		collidedEntity, exists := Game.ColladibleObjects[nextCoords]
		if exists && collidedEntity.ID != 0 && collidedEntity.Type != 0 {
			switch collidedEntity.Type {
			case c.WALL, c.BOMB, c.BLAST:
				return blastDirection
			case c.CRATE:
				Game.DestroyCrate(collidedEntity.ID)
				return blastDirection
			case c.POWERUP_SPEED, c.POWERUP_BOMB, c.POWERUP_HEAL, c.POWERUP_POWER:
				Game.Powerups[collidedEntity.ID].Death(Game.Updates)
				blastDirection = append(blastDirection, nextCoords)
				Game.ColladibleObjects[nextCoords] = ColladibleObject{ID: blast.ID, Type: c.BLAST}
				return blastDirection
			}
		}
		for _, player := range Game.Players {
			if !player.IsIFrame {
				if c.Round(math.Floor(player.Position.Y)) == nextCoords.X && c.Round(math.Floor(player.Position.Y)) == nextCoords.Y {
					fmt.Println(player.Position, nextCoords)
					player.Hit(Game.Updates)
				}
			}
		}
		blastDirection = append(blastDirection, nextCoords)
		Game.ColladibleObjects[nextCoords] = ColladibleObject{ID: blast.ID, Type: c.BLAST}
	}
	return blastDirection
}

func (Game *GameState) DestroyCrate(crateID int) {
	crate := Game.Crates[crateID]
	crate.Death(Game.Updates)
	var newPowerup Powerup
	if err := newPowerup.Initialize(crate.Position, Game.Updates); err != nil {
		Game.ColladibleObjects[crate.Position] = ColladibleObject{}
		return
	}
	delete(Game.Crates, crateID)
	Game.Powerups[newPowerup.ID] = &newPowerup
	Game.ColladibleObjects[crate.Position] = ColladibleObject{ID: newPowerup.ID, Type: newPowerup.Type}
}

func (Game *GameState) DestroyPowerup(powerupID int) {
	powerup := Game.Powerups[powerupID]
	powerup.Death(Game.Updates)
	delete(Game.Crates, powerupID)
	Game.ColladibleObjects[powerup.Position] = ColladibleObject{}
}
