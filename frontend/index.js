import Gachi, { useState, useEffect } from "./Gachi.js/src/core/framework.ts";
import { Router, Route } from "./Gachi.js/src/components/router.ts";
import Login from "./Pages/login/LoginPage.jsx";
import Lobby from "./Pages/lobby/LobbyPage.jsx";
import Game from "./Pages/game/GamePage.jsx";
import Test from "./Pages/test/Test.jsx";

function Index() {
  // const [isLogged, setIsLogged] = useState(false);
  const [userList, setUserList] = useState({});
  const [mapG, setMapG] = useState([])



  return (
    <Router
      routes={[
        {
          path: "/",
          element: <Login />,
        },
        {
          path: "/lobby",
          element: <Lobby userList={userList} setUserList={setUserList} setMap={setMapG}/>,
        },
        {
          path: "/game",
          element: <Game userList={userList} setMapG={setMapG} mapG={mapG}/>
        },
        {
          path: "/test",
          element: <Test />
        },
      ]}
    />
  );
}

Gachi.render(<Index />, document.getElementById("root"));
