package gameengine

import (
	c "bomberman/Constants"
	"math"
	"math/rand"
)

type MapData struct {
	MapMatrix Map
	Crates    map[int]*Crate
}
type Map [16][16]int

// Generate map and populate it with crate
func (MD *MapData) Initialize() {
	randomTemplate := chooseRandomTemplate()
	copy(MD.MapMatrix[:], randomTemplate[:])

	MD.Crates = make(map[int]*Crate)
	MD.populateWithCrates()
}

// Chooses random number to choose random index from MAP_TEMPLATES map[int][16][16]int
func chooseRandomTemplate() (NewMap [16][16]int) {
	randomNumber := math.Round(rand.Float64()*2) + 1
	NewMap = c.MAP_TEMPLATES[int(randomNumber)]
	return
}

func (m *MapData) populateWithCrates() {
	// crateSpawnChance is calculated based on a 16x16 matrix
	cratepawnChance := (c.CRATE_SPAWN_CHANCE * 16 * 16) / 100

	// Loop to randomly place crates in the matrix until the desired crate spawn chance is met
	for count := 0; count < cratepawnChance; {
		// Randomly select a row and column
		row := rand.Intn(16)
		col := rand.Intn(16)
		Position := c.Coords{X: float64(col), Y: float64(row)}
		pos := m.MapMatrix[row][col]
		if pos != c.BACKGROUND_1 && pos != c.BACKGROUND_2 {
			continue
		}

		// Generate the crate type
		var newCrate Crate
		newCrate.Initialize(Position)

		m.MapMatrix[row][col] = newCrate.ID

		m.Crates[newCrate.ID] = &newCrate

		count++
	}
}

func (m *Map) SpawnPlayer(PlayerID int) c.Coords {
	spawnPoint := c.PlayerSpawn[PlayerID]
	x := int(spawnPoint.X)
	y := int(spawnPoint.Y)
	m[y][x] = PlayerID
	return c.Coords{X: float64(x), Y: float64(y)}
}
