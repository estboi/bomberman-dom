import Gachi, { useNavigate, useState } from "../../Gachi.js/src/core/framework";
import "./header.css"
import Logo from  "../../Public/assets/Favicon.png"

export default function Header() {
    return (
        <div className="header">
            <div className="header-logo-desc">b</div>
            <img className="header-logo" src={Logo} alt="logo"></img>
            <div className="header-logo-desc">mberman</div>
        </div>
    )
}