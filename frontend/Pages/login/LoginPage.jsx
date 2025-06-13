import Gachi, {  useNavigate, useState, useEffect  } from "../../Gachi.js/src/core/framework";
import "./login.css"
import Header from "../components/Header";
import { ws } from "../../utils/websocket";

export default function Login() {
    const navigate = useNavigate()
    const [name, setName] = useState("")

    const handleSubmit = (e) => {
        e.preventDefault()
        ws.send(JSON.stringify({
            type: "Login",
            payload: name,
          })
        )
        navigate('/lobby')
    }

    const handleChange = (e) => {
        setName(e.target.value)
    }

    return (
        <main>
            <Header />
            <div className="login-container">

                <form className="login-form" onSubmit={handleSubmit} >
                    <div className="login-input-holder">
                        <input
                            className="login-name-input"
                            type="text"
                            value={name}
                            onChange={handleChange} />

                        <div className="login-input-placeholder">Enter your nickname</div>
                    </div>

                    <button className="login-name-button">Continue</button>
                </form>

            </div>
        </main>
    )
}