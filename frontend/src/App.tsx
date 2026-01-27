import { useEffect, useState } from "react";
import { Player } from "./component/player";
import { Register } from "./component/register";
import { GetVersion } from "../wailsjs/go/main/App";

function App() {
  const [version, setVersion] = useState("");

  useEffect(() => {
    // Gọi hàm từ Go Backend khi ứng dụng khởi động
    GetVersion()
      .then((v: number) => {
        setVersion(v.toString());
      })
      .catch((err: any) => {
        console.error("Lỗi lấy version:", err);
      });
  }, []);

  return (
    <div
      id="App"
      style={{
        width: "100vw",
        height: "100vh",
        position: "relative",
        overflow: "hidden",
        backgroundColor: "#000",
      }}
    >
      <Player />
      <Register />

      {/* Version hiển thị ở góc dưới bên phải */}
      <div
        style={{
          position: "absolute",
          bottom: "4px",
          right: "10px",
          color: "rgba(255, 255, 255, 1)",
          fontSize: "12px",
          fontFamily: "Arial, sans-serif",
          zIndex: 9999,
          pointerEvents: "none", // Không cản trở click vào các phần tử bên dưới
        }}
      >
        {version && `v${version}`}
      </div>
    </div>
  );
}

export default App;
