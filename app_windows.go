//go:build windows

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/registry"
)

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
		return
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return
	}
	defer key.Close()
	_ = key.SetStringValue("ForlifeMediaPlayer", exe)
}

func (a *App) silentUpdate() {
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(WINDOW_CHECK_UPDATE_URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return
	}

	if !apiRes.Success || apiRes.Data.Version <= AppVersionInt || apiRes.Data.URL == "" {
		return
	}

	tmp := filepath.Join(os.TempDir(), "advert.exe")
	out, err := os.Create(tmp)
	if err != nil {
		return
	}
	defer out.Close()

	r, err := http.Get(apiRes.Data.URL)
	if err != nil {
		return
	}
	defer r.Body.Close()

	_, err = out.ReadFrom(r.Body)
	if err != nil {
		return
	}

	exec.Command(tmp, "/S").Start()
	os.Exit(0)
}
