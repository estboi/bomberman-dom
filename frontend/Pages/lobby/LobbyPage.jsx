import Gachi, { useNavigate, useContext, useState, useEffect } from "../../Gachi.js/src/core/framework";
import "./lobby.css"
import Header from "../components/Header";
import Player from "../components/Player";
import Arrow from "../../Public/assets/Up_arrow.png"
import Chat from "../components/Chat"
import { subscribe, ws } from "../../utils/websocket";


export default function Lobby({ userList, setUserList, setMap }) {
    const navigate = useNavigate()
    const [timer, setTimer] = useState("Not enough players");
    const [lastCounter, setLastCounter] = useState(false)
    
    const avatars = ['avatar1', 'avatar2', 'avatar3', 'avatar4']

    const getPlayers = (data) => {
        setUserList(() => {
            return data
        })
    }

    const updateTimer = (data) => {
        setTimer(data)
    }

    const startGame = () => {
        navigate('/game')
    }

    const lastCount = () => {
        setLastCounter(true)
    }

    const genMap = (data) => {
        console.log(data);
        setMap(data)
    }

    useEffect(() => {
        if (ws.readyState !== 1) {
            navigate('/')
            return
        }
        subscribe("Generate_Map", genMap)
        subscribe("New_Player", getPlayers)
        subscribe("Timer", updateTimer)
        subscribe("Start_Game", startGame)
        subscribe("Last_Count", lastCount)

        ws.send(JSON.stringify({
            type: "New_Player",
            payload: "yo",
        })
        )

    }, []);


    return (
        <div className="lobby-page">
            <Header />

            {lastCounter ? (
                <p className="hinters">Ready to start?</p>
            ) : (
                <p className="hinters">Waiting for other players to join...</p>
            )}

            <div className="page-down">

                <div className="chat-players-container">
                    <div className="players-container">
                        {userList && Object.values(userList).map((name, index) => (
                            <Player key={index} name={name} avatar={avatars[index]} />
                        ))}
                    </div>

                    <Chat />

                </div>

                <div className="timer-controls-container">
                    <div className="timer-holder">
                        <div className="timer-placeholder">Game starts in</div>
                        <div className="timer-container">{timer}</div>
                    </div>

                    <div className="controls-container">
                        <div className="letters">
                            <div className="box w">W</div>

                            <div className="down-holder">
                                <div className="box a">A</div>
                                <div className="box s">S</div>
                                <div className="box d">D</div>
                            </div>

                            <div className="move-desc">Move</div>
                        </div>

                        <div className="arrows">
                            <div className="box w"><img src={Arrow} alt="arrow" /></div>

                            <div className="down-holder">
                                <div className="box a"><img className="arrowleft" src={Arrow} alt="arrow" /></div>
                                <div className="box s"><img className="arrowdown" src={Arrow} alt="arrow" /></div>
                                <div className="box d"><img className="arrowright" src={Arrow} alt="arrow" /></div>
                            </div>

                            <div className="move-desc">Move</div>
                        </div>

                        <div className="space">
                            <div className="space-holder">Space</div>
                            <div className="move-desc">Place bomb</div>
                        </div>
                    </div>

                </div>

            </div>

        </div>
    )
}