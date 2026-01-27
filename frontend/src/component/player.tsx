import { useEffect, useState, useRef } from "react";
import { useAppStore } from "./store";
import { Volume2, Play, Pause, Music, Info } from "lucide-react";

export function Player() {
  const targetAdvert = useAppStore((state) => state.targetAdvert);
  const advertStatus = useAppStore((state) => state.advertStatus);
  const deviceFields = useAppStore((state) => state.deviceFields);
  const updateAdvertStatus = useAppStore((state) => state.updateAdvertStatus);
  const setModalOpen = useAppStore((state) => state.setModalOpen);

  const [localUrl, setLocalUrl] = useState<string | null>(null);
  const audioRef = useRef<HTMLAudioElement>(null);

  const storeName =
    deviceFields.flatMap((g) => g.items).find((i) => i.label === "Cửa hàng")
      ?.value || "Hệ thống";

  useEffect(() => {
    if (targetAdvert?.id) {
      updateAdvertStatus({ id: targetAdvert.id });
    }
  }, [targetAdvert]);

  useEffect(() => {
    if (!targetAdvert?.audio_url) return;
    fetch(targetAdvert.audio_url)
      .then((r) => r.blob())
      .then((blob) => {
        if (localUrl) URL.revokeObjectURL(localUrl);
        setLocalUrl(URL.createObjectURL(blob));
      });
  }, [targetAdvert?.audio_url]);

  useEffect(() => {
    if (audioRef.current && advertStatus?.volume !== undefined) {
      audioRef.current.volume = advertStatus.volume;
    }
  }, [localUrl]);

  const togglePlay = () => {
    if (!audioRef.current) {
      console.log("Current audio ref null");
      return;
    }
    if (advertStatus?.playing) {
      console.log("trigger pause audio");
      audioRef.current.pause();
    } else {
      console.log("trigger play audio");
      audioRef.current.play();
    }
  };

  console.log("Advert status: ", advertStatus);

  const Footer = () => (
    <div
      style={{
        width: "100%",
        borderTop: "1px solid rgba(255,255,255,0.2)",
        paddingTop: "15px",
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
      }}
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: "6px",
          fontSize: "11px",
        }}
      >
        <Music size={14} color="white" />{" "}
        <span style={{ opacity: 0.9 }}>
          {advertStatus?.playing ? "Đang phát nhạc" : "Đã tạm dừng"}
        </span>
      </div>

      <button
        onClick={() => setModalOpen(true)}
        style={{
          background: "transparent",
          border: "1px solid white",
          color: "white",
          padding: "6px 14px",
          borderRadius: "6px",
          fontSize: "10px",
          cursor: "pointer",
          display: "flex",
          alignItems: "center",
          gap: "6px",
          fontWeight: "bold",
          letterSpacing: "0.5px",
        }}
      >
        <Info size={12} /> XEM THÔNG TIN
      </button>
    </div>
  );

  if (!targetAdvert) {
    return (
      <div
        style={{
          width: "100%",
          height: "100%",
          background: "linear-gradient(180deg, #94C9FF 0%, #C6A2FF 100%)",
          display: "flex",
          flexDirection: "column",
          padding: "20px",
          boxSizing: "border-box",
          color: "white",
          fontFamily: "sans-serif",
          overflow: "hidden",
        }}
      >
        <div
          style={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <div
            style={{ fontSize: "18px", fontWeight: "bold", marginBottom: 6 }}
          >
            Chưa có quảng cáo nào
          </div>
          <div style={{ fontSize: "12px", opacity: 0.9 }}>
            Tại "{storeName}"
          </div>
        </div>
        <Footer />
      </div>
    );
  }

  return (
    <div
      style={{
        width: "100%",
        height: "100%",
        background: "linear-gradient(180deg, #94C9FF 0%, #C6A2FF 100%)",
        display: "flex",
        flexDirection: "column",
        padding: "20px",
        boxSizing: "border-box",
        color: "white",
        fontFamily: "sans-serif",
        overflow: "hidden",
      }}
    >
      <div style={{ textAlign: "center", marginTop: "10px" }}>
        <div style={{ fontSize: "13px", opacity: 0.8, fontWeight: 500 }}>
          Đang phát quảng cáo
        </div>
        <div style={{ fontSize: "18px", fontWeight: "bold", margin: "6px 0" }}>
          {targetAdvert.name}
        </div>
        <div style={{ fontSize: "12px", opacity: 0.9 }}>Tại "{storeName}"</div>
      </div>

      <div
        style={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          width: "100%",
        }}
      >
        <div
          onClick={togglePlay}
          style={{
            cursor: "pointer",
            marginBottom: "40px",
            transition: "transform 0.1s",
          }}
          onMouseDown={(e) => (e.currentTarget.style.transform = "scale(0.9)")}
          onMouseUp={(e) => (e.currentTarget.style.transform = "scale(1)")}
        >
          {advertStatus?.playing ? (
            <Pause size={80} strokeWidth={1.5} fill="white" />
          ) : (
            <Play size={80} strokeWidth={1.5} fill="white" />
          )}
        </div>

        <div
          style={{
            width: "95%",
            display: "flex",
            alignItems: "center",
            gap: "12px",
          }}
        >
          <Volume2 size={18} color="white" style={{ opacity: 0.8 }} />
          <div
            style={{
              flex: 1,
              position: "relative",
              display: "flex",
              alignItems: "center",
            }}
          >
            <input
              type="range"
              min="0"
              max="1"
              step="0.01"
              value={advertStatus?.volume || 0}
              onChange={(e) => {
                const v = parseFloat(e.target.value);
                updateAdvertStatus({ volume: v });
                if (audioRef.current) audioRef.current.volume = v;
              }}
              style={{
                WebkitAppearance: "none",
                width: "100%",
                height: "4px",
                borderRadius: "2px",
                background: `linear-gradient(to right, #70D1FF ${
                  (advertStatus?.volume || 0) * 100
                }%, rgba(255,255,255,0.3) ${
                  (advertStatus?.volume || 0) * 100
                }%)`,
                outline: "none",
                cursor: "pointer",
              }}
            />
            <style>{`
              input[type='range']::-webkit-slider-thumb {
                -webkit-appearance: none;
                height: 14px;
                width: 14px;
                border-radius: 50%;
                background: white;
                box-shadow: 0 0 5px rgba(0,0,0,0.2);
              }
            `}</style>
          </div>
          <span
            style={{
              fontSize: "14px",
              width: "25px",
              textAlign: "right",
              fontWeight: "bold",
            }}
          >
            {Math.round((advertStatus?.volume || 0) * 100)}
          </span>
        </div>
      </div>

      <audio
        ref={audioRef}
        src={localUrl || ""}
        autoPlay
        loop
        onPlay={() => updateAdvertStatus({ playing: true })}
        onPause={() => updateAdvertStatus({ playing: false })}
        onVolumeChange={(e) =>
          updateAdvertStatus({ volume: e.currentTarget.volume })
        }
      />

      <Footer />
    </div>
  );
}
