package websocket

import (
	constants "bomberman/Constants"
	gameengine "bomberman/GameEngine"
	"encoding/json"
	"errors"
	"log"
	"sort"
	"time"
)

var fullLobby bool
var PlayerInputChannel chan constants.PlayerInput

// SortClients sorts the clients by name
func SortClients(clients map[*Client]bool) []Client {
	var sortedClients []Client
	for client := range clients {
		sortedClients = append(sortedClients, *client)
	}

	// Define a custom sorting function based on JoinTime
	sort.Slice(sortedClients, func(i, j int) bool {
		return sortedClients[i].JoinTime.Before(sortedClients[j].JoinTime)
	})

	return sortedClients
}

func Login(event Event, c *Client) error {
	var loginDTO string
	if err := json.Unmarshal(event.Payload, &loginDTO); err != nil {
		log.Println("Error unmarshaling event: ", err)
	}
	c.Name = loginDTO
	c.JoinTime = time.Now()
	tempCount := 0
	for Client := range c.Manager.Clients {
		if Client.UserID > 10 {
			tempCount++
		}
	}
	c.UserID = tempCount + 11
	for client := range c.Manager.Clients {
		client.MessageChan <- event
	}
	return nil
}

func NewPlayer(event Event, c *Client) error {
	if len(c.Manager.Clients) > 4 {
		return errors.New("lobby is full")
	}

	sortedClients := SortClients(c.Manager.Clients)

	var userList []string
	for _, client := range sortedClients {
		if client.UserID > 10 {
			userList = append(userList, client.Name)
		}
	}

	dataBytes, err := json.Marshal(userList)
	if err != nil {
		log.Printf("error marshaling data: %s\n", err)
		return err
	}

	event.Payload = json.RawMessage(dataBytes)
	for client := range c.Manager.Clients {
		client.MessageChan <- event
	}
	go Timer(len(userList), c)
	return nil
}

type MessageDTO struct {
	UserID  int    `json:"id"`
	Message string `json:"message"`
}

func NewMessage(event Event, c *Client) error {
	var message MessageDTO

	if err := json.Unmarshal(event.Payload, &message); err != nil {
		log.Println("Error unmarshaling event: ", err)
	}

	message.UserID = c.UserID

	dataBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling data: %s\n", err)
		return err
	}
	event.Type = "New_Message"
	event.Payload = json.RawMessage(dataBytes)
	for client := range c.Manager.Clients {
		client.MessageChan <- event
	}
	return nil
}

func Timer(playerCount int, c *Client) error {
	//Initial timer setup
	var timerEvent Event
	var duration int
	switch playerCount {
	// case 0, 1, 3:
	case 0:
		return nil
	case 2: // start timer
		duration = 3
	case 4: // final timer
		duration = 1
	}

	timerEvent.Type = "Timer"
	durationJSON, _ := json.Marshal(duration)
	timerEvent.Payload = json.RawMessage(durationJSON)
	for client := range c.Manager.Clients {
		client.MessageChan <- timerEvent
	}

	for i := duration; i >= 0; i-- {
		if fullLobby {
			break
		}
		if i == 0 && !fullLobby {
			startLastTimer(timerEvent, c)
		}
		if playerCount == 4 {
			fullLobby = true
		}
		durationJSON, _ := json.Marshal(i)
		timerEvent.Payload = json.RawMessage(durationJSON)
		for client := range c.Manager.Clients {
			client.MessageChan <- timerEvent
		}
		time.Sleep(time.Second)
	}

	if fullLobby {
		startLastTimer(timerEvent, c)
	}

	return nil
}

func startLastTimer(timerEvent Event, c *Client) error {
	timerEvent.Type = "Last_Count"
	for client := range c.Manager.Clients {
		client.MessageChan <- timerEvent
	}

	PlayerInputChannel = make(chan constants.PlayerInput)
	var GameState gameengine.GameState
	GameState.Initialize()

	for j := 2; j >= 0; j-- {
		timerEvent.Type = "Timer"
		durationJSON, _ := json.Marshal(j)
		timerEvent.Payload = json.RawMessage(durationJSON)
		for client := range c.Manager.Clients {
			client.MessageChan <- timerEvent
		}
		if j == 0 {
			go StartGame(&GameState, PlayerInputChannel, c.Manager)

			SendGameState(&GameState, c.Manager)

			SendStartGame(timerEvent, c)
		}
		time.Sleep(time.Second)
	}
	return nil
}

func SendStartGame(event Event, c *Client) error {
	event.Type = "Start_Game"
	for client := range c.Manager.Clients {
		client.MessageChan <- event
	}
	return nil
}

func InputHandle(event Event, c *Client) error {
	// Recieve input
	var input string
	if err := json.Unmarshal(event.Payload, &input); err != nil {
		log.Println("Error unmarshaling event: ", err)
	}

	var playerInput = constants.PlayerInput{
		Direction: input,
		PlayerID:  c.UserID,
	}

	PlayerInputChannel <- playerInput

	return nil
}
