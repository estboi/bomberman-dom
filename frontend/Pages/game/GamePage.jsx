import Gachi, { useState, useEffect, useNavigate } from "../../Gachi.js/src/core/framework";
import "./game.css";
import Background from "../../Public/assets/Background.png"
import Header from "../components/Header";
import Icons from "../components/Icons";
import { subscribe, ws } from "../../utils/websocket";


export default function Game({ userList, setMapG, mapG }) {
    const navigate = useNavigate()

    const widthConv = 2.5
    const [map, setMap] = useState(mapG);
    const [players, setPlayers] = useState([
        {
            top: "2.5rem",
            left: "2.5rem"
        },
        {
            top: "2.5rem",
            left: "35rem"
        },
        {
            top: "35rem",
            left: "2.5rem"
        },
        {
            top: "35rem",
            left: "35rem"
        }
    ])


    useEffect(() => {
        setPlayers(players.slice(0, userList.length))
    }, [userList]);

    const avatars = ['avatar1', 'avatar2', 'avatar3', 'avatar4']

    useEffect(() => {
        const handleKeyDown = (e) => {
            switch (e.key) {
                case "w":
                case "ArrowUp":
                    ws.send(JSON.stringify({ "type": "Input", "payload": "U" }))
                    break;
                case "a":
                case "ArrowLeft":
                    ws.send(JSON.stringify({ "type": "Input", "payload": "L" }))
                    break;
                case "s":
                case "ArrowDown":
                    ws.send(JSON.stringify({ "type": "Input", "payload": "D" }))
                    break;
                case "d":
                case "ArrowRight":
                    ws.send(JSON.stringify({ "type": "Input", "payload": "R" }))
                    break;
                case " ":
                    ws.send(JSON.stringify({ "type": "Input", "payload": "P" }))
                    break;
                default:
                    break;
            }
        }
        document.addEventListener("keydown", handleKeyDown);
        return () => {
            document.removeEventListener("keydown", handleKeyDown);
        };
    }, [])

    const bombManager = (data) => {
        let bomb = document.createElement("div")
        bomb.className = "entity_tileset bomb bomb_activation"
        bomb.style.top = `${data.Coords.Y * widthConv}rem`
        bomb.style.left = `${data.Coords.X * widthConv}rem`
        bomb.id = `${data.ID}bomb`
        let parentEl = document.getElementById("game-container")
        parentEl.appendChild(bomb)
    }

    const blastManager = (data) => {
        let bomb = document.getElementById(`${data.ID}bomb`)
        console.log(bomb)
        bomb.remove()
        handleBlast(data.Coords.Center, "center", data.ID)
        handleBlast(data.Coords.Up, "up", data.ID)
        handleBlast(data.Coords.Right, "right", data.ID)
        handleBlast(data.Coords.Down, "down", data.ID)
        handleBlast(data.Coords.Left, "left", data.ID)
    }

    const handleBlast = (arr, type, id) => {
        const classMap = {
            right: { end: "rotate_90", middle: "rotate_90" },
            down: { end: "rotate_180" },
            left: { end: "rotate_270", middle: "rotate_90" }
        };
        if (type == "center") {
            let blast = document.createElement("div")
            blast.className = "entity_tileset blast_center"
            blast.style.top = `${arr.Y * widthConv}rem`
            blast.style.left = `${arr.X * widthConv}rem`
            blast.id = `${id}blast`
            let parentEl = document.getElementById("game-container")
            parentEl.appendChild(blast)
            return
        }

        for (let i = 0; i < arr.length; i++) {
            let blast = document.createElement("div")
            blast.style.top = `${arr[i].Y * widthConv}rem`
            blast.style.left = `${arr[i].X * widthConv}rem`
            blast.id = `${id}blast`
            let className = i === arr.length - 1 ? "entity_tileset blast_end" : "entity_tileset blast_middle";

            if (classMap[type]) {
                let additionalClass = classMap[type][className.includes("end") ? "end" : "middle"];
                if (additionalClass) {
                    className += ` ${additionalClass}`;
                }
            }
            blast.className = className;
            let parentEl = document.getElementById("game-container")
            parentEl.appendChild(blast)
        }
    };

    const blastDeadManager = (data) => {
        console.log("blast dead:", data);
        let totalCount = 1;
        for (const direction in data.Coords) {
            if (Array.isArray(data.Coords[direction])) {
                totalCount += data.Coords[direction].length;
            }
        }
        for (let index = 0; index < totalCount; index++) {
            let blast = document.getElementById(`${data.ID}blast`)
            blast.remove()
        }
    }

    const crateDeadManager = (data) => {
        console.log("crate dead:", data);
        let crate = document.getElementById(data)
        crate.firstChild.classList.remove("crate")
        crate.firstChild.classList.add("crate_broken")
        setTimeout(() => {
            crate.firstChild.classList.remove("crate_broken");
        }, 5000);

        for (let i = 0; i < map.length; i++) {
            for (let j = 0; j < map[i].length; j++) {
                if (map[i][j] === data) {
                    map[i][j] = 0;
                    break;
                }
            }
        }
    }

    const powerupManager = (data) => {
        console.log("powerup:", data);
        let powerup = document.createElement("div")
        switch (data.Type) {
            case 21:
                powerup.className = "entity_tileset powerup_health"
                break;
            case 22:
                powerup.className = "entity_tileset powerup_blast"
                break;
            case 23:
                powerup.className = "entity_tileset powerup_speed"
                break;
            case 24:
                powerup.className = "entity_tileset powerup_bomb"
                break;
        }
        powerup.style.top = `${data.Coord.Y * widthConv}rem`
        powerup.style.left = `${data.Coord.X * widthConv}rem`
        powerup.id = `${data.ID}powerup`
        let parentEl = document.getElementById("game-container")
        parentEl.appendChild(powerup)
    }

    const powerupDeadManager = (data) => {
        console.log("powerup dead:", data);
        let powerup = document.getElementById(`${data}powerup`)
        powerup.remove()
    }

    const playerDeadManager = (data) => {
        console.log("player dead:", data);
        let player = document.getElementById(data)
        player.className = "entity_tileset dead_body"
        setTimeout(() => {
            player.remove()
        }, 5000);
    }

    const moveManager = (data) => {
        setPlayers(players => {
            const updatedPlayers = [...players];
            updatedPlayers[data.PlayerID - 11] = {
                ...updatedPlayers[data.PlayerID - 11],
                top: `${data.Coords.Y * widthConv}rem`,
                left: `${data.Coords.X * widthConv}rem`
            };
            return updatedPlayers;
        });
        let player = document.getElementById(data.PlayerID)
        // player.className += " moving"
        player.className = player.className.replace(/\brotate_\d{1,3}\b/g, '');
        switch (data.Direction) {
            case "U": player.className += " rotate_270"
                break;
            case "R": player.className += " rotate_0"
                break;
            case "D": player.className += " rotate_90"
                break;
            case "L": player.className += " rotate_180"
                break;
        }

    }

    const endGameManager = (data) => {
        console.log("end:", data);
        alert("end game")
        
    }

    useEffect(() => {
        if (ws.readyState !== 1) {
            navigate('/')
            return
        }
        subscribe("move", moveManager)
        subscribe("bomb", bombManager)
        subscribe("blast", blastManager)
        subscribe("blast_dead", blastDeadManager)
        subscribe("crate_dead", crateDeadManager)
        subscribe("powerup", powerupManager)
        subscribe("powerup_dead", powerupDeadManager)
        subscribe("player_dead", playerDeadManager)
        subscribe("end_game", endGameManager)
    }, [])


    const renderTexture = (block) => {
        switch (block) {
            case 0:
            case 3:
            case 11:
            case 12:
            case 13:
            case 14:
                return <div className="map_tileset bg1"></div>;
            case 1:
                return <div className="map_tileset bg2"></div>;
            case 2:
                return <div className="map_tileset wall"></div>;
            case 15:
                return <div className="entity_tileset dead_body"></div>;
            case 41:
                return <div className="entity_tileset crate_broken"></div>;
            default:
                return <div className="entity_tileset crate"></div>;
        }
    };




    return (
        <main>
            <div className="hid">
                <Header />
            </div>

            <div className="game-page">
                <div className="main-game-container">
                    {/* <div className="avatars-container">
                        {players.map((_, index) => (
                            <Icons key={index} avatar={avatars[index]} />
                        ))}
                    </div>/ */}

                    <div className="game-field" id="game-field">
                        <div className="game-container" id="game-container" >
                            <img id="background" src={Background} />
                            {players && (players.map((_, index) => {
                                return (
                                    <div
                                        style={`top: ${players[index].top}; left: ${players[index].left}`}
                                        className={`entity_tileset player${index + 1}`}
                                        id={`1${index + 1}`}>
                                    </div>
                                );
                            })
                            )}
                            {map.map((item, rowIndex) => (
                                <div key={rowIndex} className="row">
                                    {item.map((block, columnIndex) => (
                                        <div
                                            key={columnIndex}
                                            className="tile"
                                            id={rowIndex * 100 + columnIndex}>
                                            {renderTexture(block)}
                                        </div>
                                    ))}
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* <div className="stats-container">
                        <div className="stats-placeholder">Stats</div>
                        <div className="icons-container">
                            <div className="c4-num-icon">
                                <div className="icons-tileset c4-icon"></div>
                                <div className="number">x1</div>
                            </div>
                            <div className="blast-num-icon">
                                <div className="icons-tileset blast-icon"></div>
                                <div className="number">x1</div>
                            </div>
                            <div className="speed-num-icon">
                                <div className="icons-tileset speed-icon"></div>
                                <div className="number">x1</div>
                            </div>
                        </div>
                    </div> */}
                </div>
            </div>
        </main>
    );
}