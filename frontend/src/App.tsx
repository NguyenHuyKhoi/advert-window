import { Player } from "./component/player";
import { Register } from "./component/register";

function App() {
  return (
    <div
      id="App"
      style={{
        width: "100vw",
        height: "100vh",
        position: "relative",
        overflow: "hidden",
      }}
    >
      <Player />
      <Register />
    </div>
  );
}

export default App;
