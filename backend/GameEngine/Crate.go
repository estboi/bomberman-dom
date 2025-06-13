package gameengine

import (
	constants "bomberman/Constants"
	"encoding/json"
)

type Crate struct {
	ID       int
	Type     int
	Position constants.Coords

	IsAlive bool
}

func (C *Crate) Initialize(Position constants.Coords) {
	C.Type = constants.CRATE
	C.Position = Position
	C.IsAlive = true
	C.CreateID()
}

func (C *Crate) CreateID() {
	posSum := C.Position.Y*100 + C.Position.X
	C.ID = int(posSum)
}

func (C *Crate) Death(update chan UpdatedEntity) {
	crateIdJSON, _ := json.Marshal(C.ID)
	update <- UpdatedEntity{Type: "crate_dead", Payload: crateIdJSON}
}
