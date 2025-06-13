import Gachi, { useState, useEffect, useContext } from "../../Gachi.js/src/core/framework";
import "./chat.css"
import { subscribe, ws } from "../../utils/websocket";

export default function Chat() {
  const [message, setMessage] = useState({id: null, message: null});
  const [messages, setMessages] = useState([]);

  const handleMessageChange = (e) => {
    setMessage((prevMessage) => ({
      ...prevMessage,
      message: e.target.value
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    if (message.message.trim() !== "") {
      ws.send(JSON.stringify({
        type: "New_Message",
        payload: {message: message.message},
      }));
      event.target.children[0].value = ""
    }
  };

  const getMessages = (data) => {
    console.log(data)
    setMessages((prev) => {
      return [...prev, data]
    });
  };

  useEffect(() => {
    subscribe("New_Message", getMessages)
  }, []);

  return (
    <div className="chat-container">
      <div className="chat-holder">
        <div className="scroller">
          <div className="chat-placeholder">Chat</div>
          <div className="messages-container">
            {messages.slice().reverse().map((msg) => (
              <div
                key={msg.id}
                className={`message-player-${msg.id}`}>
                {msg.message}
              </div>
            ))}
          </div>

        </div>
      </div>

      <form className="chat-form" onSubmit={handleSubmit}>
        <input
          className="chat-input"
          type="text"
          placeholder="Type here..."
          onChange={handleMessageChange}
        />
        <button className="chat-submit-message" type="submit">
          Send
        </button>
      </form>
    </div>
  );
}