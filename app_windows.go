//go:build windows

package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("[AutoStart] exe error:", err)
		return
	}
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		fmt.Println("[AutoStart] registry error:", err)
		return
	}
	defer key.Close()
	_ = key.SetStringValue("ForlifeMediaPlayer", exe)
}

func (a *App) silentUpdate() {
	fmt.Println("[Update] Checking update...")
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(WINDOW_CHECK_UPDATE_URL)
	if err != nil {
		fmt.Println("[Update] API error:", err)
		return
	}
	defer resp.Body.Close()

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		fmt.Println("[Update] JSON error:", err)
		return
	}

	fmt.Printf("[Update] ServerVersion=%d LocalVersion=%d Success=%v URL=%s\n",
		apiRes.Data.Version, AppVersionInt, apiRes.Success, apiRes.Data.URL)

	if !apiRes.Success || apiRes.Data.Version <= AppVersionInt || apiRes.Data.URL == "" {
		fmt.Println("[Update] No update needed")
		return
	}

	tmp := filepath.Join(os.TempDir(), "advert.exe")
	fmt.Println("[Update] Downloading to:", tmp)

	out, err := os.Create(tmp)
	if err != nil {
		fmt.Println("[Update] Create file error:", err)
		return
	}
	defer out.Close()

	r, err := http.Get(apiRes.Data.URL)
	if err != nil {
		fmt.Println("[Update] Download error:", err)
		return
	}
	defer r.Body.Close()

	_, err = out.ReadFrom(r.Body)
	if err != nil {
		fmt.Println("[Update] Write error:", err)
		return
	}

	fmt.Println("[Update] Running installer silent...")
	exec.Command(tmp, "/S").Start()
	os.Exit(0)
}
