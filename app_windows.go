//go:build windows

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/registry"
)

var logger *log.Logger

func init() {
	dir, _ := os.UserConfigDir()
	logPath := filepath.Join(dir, "ForlifeMediaPlayer", "app.log")
	_ = os.MkdirAll(filepath.Dir(logPath), 0700)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		logger = log.New(io.MultiWriter(f, os.Stdout), "", log.LstdFlags)
		logger.Println("==== Logger started ====")
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}
}

type APIResponse struct {
	Success bool    `json:"success"`
	Status  int     `json:"status"`
	Data    AppData `json:"data"`
}

type AppData struct {
	Version int    `json:"version"`
	URL     string `json:"url"`
}

func (a *App) enableAutoStart() {
	exe, err := os.Executable()
	if err != nil {
		logger.Println("[AutoStart] exe error:", err)
		return
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		logger.Println("[AutoStart] registry error:", err)
		return
	}
	defer key.Close()

	_ = key.SetStringValue("ForlifeMediaPlayer", exe)
	logger.Println("[AutoStart] Registry updated:", exe)
}

func (a *App) silentUpdate() {
	logger.Println("[Update] Checking update...")

	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(WINDOW_CHECK_UPDATE_URL)
	if err != nil {
		logger.Println("[Update] API error:", err)
		return
	}
	defer resp.Body.Close()

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		logger.Println("[Update] JSON error:", err)
		return
	}

	logger.Printf("[Update] Server=%d Local=%d Success=%v URL=%s\n",
		apiRes.Data.Version,
		AppVersionInt,
		apiRes.Success,
		apiRes.Data.URL,
	)

	if !apiRes.Success || apiRes.Data.Version <= AppVersionInt || apiRes.Data.URL == "" {
		logger.Println("[Update] No update needed")
		return
	}

	tmp := filepath.Join(os.TempDir(), "advert.exe")
	logger.Println("[Update] Downloading installer to:", tmp)

	out, err := os.Create(tmp)
	if err != nil {
		logger.Println("[Update] Create file error:", err)
		return
	}

	r, err := http.Get(apiRes.Data.URL)
	if err != nil {
		logger.Println("[Update] Download error:", err)
		out.Close()
		return
	}

	_, err = out.ReadFrom(r.Body)
	r.Body.Close()
	out.Close()

	if err != nil {
		logger.Println("[Update] Write error:", err)
		return
	}

	logger.Println("[Update] Launching installer silent...")

	cmd := exec.Command(tmp, "/S")
	err = cmd.Start()
	if err != nil {
		logger.Println("[Update] Installer start error:", err)
		return
	}

	logger.Println("[Update] Installer started. Exiting old app.")
	os.Exit(0)
}
