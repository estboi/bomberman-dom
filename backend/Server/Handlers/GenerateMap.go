package handlers

import (
	gameengine "bomberman/GameEngine"
	"encoding/json"
	"net/http"
)

func GenerateMap(w http.ResponseWriter, r *http.Request) {
	var MapData gameengine.MapData
	MapData.Initialize()

	json, err := json.Marshal(MapData.MapMatrix)
	if err != nil {
		http.Error(w, "Couldn`t create a JSON for map", http.StatusBadRequest)
	}

	w.Write(json)

}
