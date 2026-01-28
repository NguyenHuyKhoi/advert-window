import { useEffect, useState, useRef } from "react";
import { GetDeviceInfo } from "../../wailsjs/go/main/App";
import { BASE_URL } from "../constants";
import { useAppStore } from "./store";

const INTERVAL_MS = 12000;
const INTERVAL_SEC = INTERVAL_MS / 1000;

export function Register() {
  const isOpen = useAppStore((state) => state.isInfoModalOpen);
  const setModalOpen = useAppStore((state) => state.setModalOpen);
  const setDeviceInfo = useAppStore((state) => state.setDeviceInfo);
  const setTargetAdvert = useAppStore((state) => state.setTargetAdvert);
  const setDeviceFields = useAppStore((state) => state.setDeviceFields);
  const groups = useAppStore((state) => state.deviceFields);
  const advertStatus = useAppStore((state) => state.advertStatus);

  const advertStatusRef = useRef(advertStatus);
  const [countdown, setCountdown] = useState(INTERVAL_SEC);

  useEffect(() => {
    advertStatusRef.current = advertStatus;
  }, [advertStatus]);

  const register = async () => {
    try {
      const info = await GetDeviceInfo();
      setDeviceInfo(info);

      const current_advert = {
        id: advertStatusRef?.current?.id,
        playing: advertStatusRef?.current?.playing || false,
        volume: Math.round((advertStatusRef?.current?.volume ?? 0) * 100),
      };

      const res = await fetch(
        `${BASE_URL}/devices/register/window?ip=42.112.114.188`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "ngrok-skip-browser-warning": "true",
          },
          body: JSON.stringify({
            device_info: info,
            current_advert: current_advert,
          }),
        },
      );

      const json = await res.json();

      if (json.data?.device_fields) setDeviceFields(json.data.device_fields);

      const newTarget = json.data?.target_advert;

      setTargetAdvert(newTarget);

      setCountdown(INTERVAL_SEC);
    } catch (e) {
      console.error("Register Error:", e);
      setCountdown(INTERVAL_SEC);
    }
  };

  useEffect(() => {
    // Gọi lần đầu khi mount
    register();

    const timer = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          register();
          return INTERVAL_SEC;
        }
        return prev - 1;
      });
    }, 1000);

    // Dọn dẹp timer khi unmount
    return () => clearInterval(timer);
  }, []); // Tuyệt đối để mảng rỗng để tránh loop/spam API

  if (!isOpen) return null;

  return (
    <div
      onClick={() => setModalOpen(false)}
      style={{
        position: "fixed",
        top: 0,
        left: 0,
        width: "100vw",
        height: "100vh",
        backgroundColor: "rgba(0,0,0,0.4)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        zIndex: 9999,
      }}
    >
      <div
        onClick={(e) => e.stopPropagation()}
        style={{
          width: "80%",
          background: "rgba(255, 255, 255, 0.6)",
          padding: "20px",
          borderRadius: "12px",
          color: "#000",
          backdropFilter: "blur(10px)",
          boxShadow: "0 4px 30px rgba(0, 0, 0, 0.1)",
        }}
      >
        <div style={{ textAlign: "center", marginBottom: "15px" }}>
          <div style={{ fontSize: "10px", opacity: 0.5, marginTop: "2px" }}>
            Đồng bộ với hệ thống sau: {countdown} giây
          </div>
        </div>

        <div
          style={{
            display: "flex",
            flexDirection: "column",
            gap: "15px",
            maxHeight: "60vh",
            overflowY: "auto",
          }}
        >
          {groups.map((g, i) => (
            <div key={i}>
              <div
                style={{
                  fontWeight: "800",
                  color: "#333",
                  marginBottom: "8px",
                  fontSize: "10px",
                  textTransform: "uppercase",
                }}
              >
                {g.title}
              </div>
              <div
                style={{ display: "flex", flexDirection: "column", gap: "6px" }}
              >
                {g.items.map((item: any, j: number) => (
                  <div
                    key={j}
                    style={{
                      display: "flex",
                      justifyContent: "space-between",
                      fontSize: "11px",
                    }}
                  >
                    <span style={{ fontWeight: "400", opacity: 0.8 }}>
                      {item.label}
                    </span>
                    <span style={{ fontWeight: "700" }}>{item.value}</span>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
