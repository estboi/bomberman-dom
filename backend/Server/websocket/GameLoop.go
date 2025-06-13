package websocket

import (
	c "bomberman/Constants"
	constants "bomberman/Constants"
	gameengine "bomberman/GameEngine"
	"encoding/json"
	"log"
	"time"
)

func StartGame(Game *gameengine.GameState, InputChan chan c.PlayerInput, Manager *Manager) {
	var endGame bool
	// Start the loop
	for {
		select {
		case playerInput := <-InputChan:
			// Handle playerInput
			Game.PlayerInput(playerInput)

		// Get the updated state part
		case UpdatedState := <-Game.Updates:
			Game.Mutex.Lock()
			var Update = Event{
				Type:    UpdatedState.Type,
				Payload: UpdatedState.Payload,
			}

			if UpdatedState.Type == "player_dead" {
				endGame = checkIfPlayerIsDead(UpdatedState.Payload, Game)
			}

			for client := range Manager.Clients {
				client.MessageChan <- Update
			}
			if endGame {
				var WINNER int
				for _, player := range Game.Players {
					WINNER = player.ID
				}
				payload, err := json.Marshal(WINNER)
				if err != nil {
					log.Println(err)
				}
				for client := range Manager.Clients {
					client.MessageChan <- Event{Type: "end_game", Payload: payload}
				}
			}

			Game.Mutex.Unlock()
		default:
			checkTimebasedEvents(Game)
			time.Sleep(33 * time.Millisecond)
		}
	}
}

func checkTimebasedEvents(game *gameengine.GameState) {
	if len(game.Bombs) != 0 {
		for _, bomb := range game.Bombs {
			if time.Since(bomb.StartTimer) >= time.Duration(constants.BOMB_TIMER)*time.Millisecond {
				game.Explode(bomb)
			}
		}
	}
	if len(game.Blast) != 0 {
		for _, blast := range game.Blast {
			if blast == nil {
				continue
			}
			if time.Since(blast.StartTimer) >= time.Duration(constants.BLAST_TIMER)*time.Millisecond {
				blast.Death(game.Updates)
				game.ClearBlast(blast)
			}
		}
	}
	for _, player := range game.Players {
		if player.IsIFrame && time.Since(player.IFrameTimer) >= time.Duration(constants.PLAYER_IFRAME)*time.Millisecond {
			log.Println("Players Iframe stop")
			player.IsIFrame = false
		}
	}
}

func checkIfPlayerIsDead(payload json.RawMessage, game *gameengine.GameState) bool {
	var playerId int
	if err := json.Unmarshal(payload, &playerId); err != nil {
		log.Println(err)
	}
	delete(game.Players, playerId)
	return len(game.Players) == 1
}

func SendGameState(Game *gameengine.GameState, Manager *Manager) {
	for client := range Manager.Clients {
		Game.InitializePlayer(client.UserID)
	}
	mapJson, _ := json.Marshal(Game.Map.MapMatrix)
	var event = Event{
		Type:    "Generate_Map",
		Payload: json.RawMessage(mapJson),
	}
	for client := range Manager.Clients {
		client.MessageChan <- event
	}

}
