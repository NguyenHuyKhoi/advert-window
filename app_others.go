//go:build !windows

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
	fmt.Println("[AutoStart] Skipped on non-Windows")
}

func (a *App) silentUpdate() {
	fmt.Printf("[Update] Checking: %s\n", WINDOW_CHECK_UPDATE_URL)

	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(WINDOW_CHECK_UPDATE_URL)
	if err != nil {
		fmt.Printf("[Update] Request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		fmt.Printf("[Update] JSON decode error: %v\n", err)
		return
	}

	fmt.Printf("[Update] Server response -> success=%v version_local=%d version_remote=%d url=%s\n", apiRes.Success, AppVersionInt, apiRes.Data.Version, apiRes.Data.URL)

	if !apiRes.Success || apiRes.Data.Version <= AppVersionInt || apiRes.Data.URL == "" {
		fmt.Println("[Update] No update needed")
		return
	}

	tmp := filepath.Join(os.TempDir(), "advert-update")
	out, err := os.Create(tmp)
	if err != nil {
		fmt.Printf("[Update] Create temp file error: %v\n", err)
		return
	}
	defer out.Close()

	r, err := http.Get(apiRes.Data.URL)
	if err != nil {
		fmt.Printf("[Update] Download error: %v\n", err)
		return
	}
	defer r.Body.Close()

	_, err = out.ReadFrom(r.Body)
	if err != nil {
		fmt.Printf("[Update] Write file error: %v\n", err)
		return
	}

	fmt.Printf("[Update] Downloaded installer to %s\n", tmp)

	exec.Command(tmp).Start()
	os.Exit(0)
}
