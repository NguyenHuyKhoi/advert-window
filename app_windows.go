//go:build windows
// +build windows

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

// Installer xử lý autostart
func (a *App) enableAutoStart() {}

func filenameFromURL(u string) string {
	base := filepath.Base(u)
	if base == "" || base == "." || base == "/" {
		return "forlife-installer.exe"
	}
	return base
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

	if resp.StatusCode != http.StatusOK {
		logger.Println("[Update] Bad status:", resp.Status)
		return
	}

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		logger.Println("[Update] JSON error:", err)
		return
	}

	if !apiRes.Success || apiRes.Data.Version <= AppVersionInt || apiRes.Data.URL == "" {
		logger.Println("[Update] No update needed")
		return
	}

	filename := filenameFromURL(apiRes.Data.URL)
	tmp := filepath.Join(os.TempDir(), filename)
	logger.Println("[Update] Downloading installer to:", tmp)

	out, err := os.Create(tmp)
	if err != nil {
		logger.Println("[Update] Create file error:", err)
		return
	}
	defer out.Close()

	r, err := http.Get(apiRes.Data.URL)
	if err != nil {
		logger.Println("[Update] Download error:", err)
		return
	}
	defer r.Body.Close()

	if _, err := io.Copy(out, r.Body); err != nil {
		logger.Println("[Update] Write error:", err)
		return
	}

	logger.Println("[Update] Launching installer silent...")
	cmd := exec.Command(tmp, "/S")
	if err := cmd.Start(); err != nil {
		logger.Println("[Update] Start error:", err)
		return
	}

	os.Exit(0)
}
