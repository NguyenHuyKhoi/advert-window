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

func (a *App) enableAutoStart() {
	// Handled by installer (HKCU\Run)
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

	tmp := filepath.Join(os.TempDir(), "advert-update.exe")
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

	if r.StatusCode != http.StatusOK {
		logger.Println("[Update] Download bad status:", r.Status)
		return
	}

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
