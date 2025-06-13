package constants

// Game speed and timers constatns
const (
	PLAYER_SPEED             float64 = 0.1
	BOMB_TIMER               int     = 3000
	BLAST_TIMER              int     = 500
	POWERUP_TIMER            int     = 5000
	PLAYER_IFRAME            int     = 3000
	SPEED_POWERUP_MULTIPLIER float64 = 0.1
)

// Game spawn chances
const (
	CRATE_SPAWN_CHANCE  int = 10
	POWERUP_DROP_CHANCE int = 60
	HEAL_DROP_CHANCE    int = 10
	BOMB_DROP_CHANCE    int = 15
	POWER_DROP_CHANCE   int = 25
	SPEED_DROP_CHANCE   int = 30
)

const HITBOX_CENTER float64 = 0.4
const PLAYER_HITBOX_SIZE float64 = 0.9

// Objects IDs
const (
	BACKGROUND_1 = 0
	BACKGROUND_2 = 1
	WALL         = 2
	SAFE_ZONE    = 3

	// Players
	PLAYER = 10
	P_1    = 11
	P_2    = 12
	P_3    = 13
	P_4    = 14
	P_DEAD = 15

	// Powerups
	POWERUP       = 20
	POWERUP_HEAL  = 21
	POWERUP_POWER = 22
	POWERUP_SPEED = 23
	POWERUP_BOMB  = 24

	// Bomb
	BOMB            = 30
	BLAST_CENTER    = 31
	BLAST_DIRECTION = 32
	BLAST_END       = 33
	BLAST           = 34

	CRATE        = 40
	CRATE_BROKEN = 41
)
