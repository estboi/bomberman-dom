import Gachi from "../../Gachi.js/src/core/framework";
import "./player.css"

export default function Player({ name, avatar }) {

    return (
        <div className="player">
            <div className="name">{name}</div>
            <div className="avatar-holder">
                <div className={`images ${avatar}`}></div>
            </div>

    
        </div>
    )
}