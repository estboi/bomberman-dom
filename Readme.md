# Bomberman-DOM Multiplayer Game

Welcome to Bomberman-DOM, a multiplayer game inspired by the classic Bomberman, built using the DOM framework. In this game, 2 to 4 players will battle against each other until only one remains victorious. The project includes a chat feature using WebSockets to enable communication between players.

## Table of Contents

- [Code Architecture](#file-structure)
- [Getting Started](#getting-started)
- [Building and Running the Game](#building-and-running-the-game)
- [Game Mechanics](#game-mechanics)
  - [Players](#players)
  - [Map](#map)
  - [Power Ups](#power-ups)
- [User Interface](#user-interface)
  - [Nickname Entry](#nickname-entry)
  - [Waiting Page](#waiting-page)
  - [Game Start](#game-start)
- [Acknowledgments](#acknowledgments)

## File Structure

```
├── assets  >>> All of sprites and icons are stored
├── backend
│   ├── Constants  >>> Constants of game such as: objects IDs, speed constants etc.
│   ├── Entities  >>> All dynamic objects in game and methods
│   ├── Gameloop  >>> Where game state updates
│   └── Server >>> Establish connection
│       └── Client >>> Websocket and sessions 
│           ├── Sessions  >>> Session creation and management
│           └── Websocket >>> Websocket connection
└── frontend
    ├── Engine >>> Client side rendering of animations, movements, map etc.
    └── Pages
        ├── Authentication
        ├── Game
        └── Lobby
```

## Getting Started

To get started with Bomberman-DOM, follow these steps:

1. Clone the repository to your local machine.
2. Open the project in your preferred code editor.
3. Set up a local server or deploy the game to a hosting platform.

## Game Mechanics

### Players

- Number of Players: 2 - 4
- Each player starts with 3 lives.
- Players are positioned in different corners of the map.
- Players can drop bombs to eliminate opponents.

### Map

- The map is fixed, and all players see the entire map.
- Two types of blocks: destructible (creates) and indestructible (walls).
- Walls are static, while creates are randomly generated on the map.

### Power Ups

- Power-ups appear when a player destroys a block.
- Bombs: Increases the number of bombs dropped at a time by 1.
- Flames: Increases bomb explosion range in four directions by 1 block.
- Speed: Increases player movement speed.

## User Interface

### Nickname Entry

- Upon opening the game, players enter a nickname for identification.

### Waiting Page

- After selecting a nickname, players are directed to a waiting page.
- A player counter increments as users join (maximum of 4 players).
- If there are more than 2 players and it doesn't reach 4 before 20 seconds, a 10-second timer starts for players to get ready.

### Game Start

- If 4 players join before 20 seconds, a 10-second timer starts, and the game begins.
- Players battle until only one remains, with 3 lives each.

## Building and Running the Game

To build and run the game, follow these steps:

1. Ensure you have a local server set up or deploy the game to a hosting platform.
2. Open the game in a modern web browser.
3. Enter a nickname and follow the on-screen instructions to join the game.
4. Enjoy the multiplayer Bomberman-DOM experience!

## Acknowledgments

Bomberman-DOM is based on the DOM framework created for the "make-your-game" and "mini-framework" projects. Special thanks to [Your Name] for the inspiration and guidance in developing this multiplayer game.
