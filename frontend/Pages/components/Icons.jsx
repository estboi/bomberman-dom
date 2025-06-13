import Gachi from "../../Gachi.js/src/core/framework";
import "./player.css"

export default function Icons({ avatar }) {

    return (
        <div className="icon-player">
            <div className="hearts-icon">
                <div className="images hearts-container"></div>
                <div className="images hearts-container"></div>
                <div className="images hearts-container"></div>
            </div>
            <div className="icon-avatar-holder">
                <div className={`images ${avatar}`}></div>
            </div>
        </div>
    )
}