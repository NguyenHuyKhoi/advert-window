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
		mw := io.MultiWriter(f, os.Stdout)
		logger = log.New(mw, "", log.LstdFlags)
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

	tmp := filepath.Join(os.TempDir(), "advert-installer.exe")
	logger.Println("[Update] Downloading to:", tmp)

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

	_, err = out.ReadFrom(r.Body)
	if err != nil {
		logger.Println("[Update] Write error:", err)
		return
	}

	logger.Println("[Update] Running installer silent...")

	cmd := exec.Command(tmp, "/S")
	err = cmd.Start()
	if err != nil {
		logger.Println("[Update] Installer launch error:", err)
		return
	}

	logger.Println("[Update] Waiting installer to finish...")
	cmd.Wait()

	time.Sleep(500 * time.Millisecond)

	exePath, _ := os.Executable()
	logger.Println("[Update] Relaunching app:", exePath)
	exec.Command(exePath).Start()

	logger.Println("[Update] Exiting old process")
	os.Exit(0)
}
