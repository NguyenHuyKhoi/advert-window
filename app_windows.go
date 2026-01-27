//go:build windows
package main

import (
	"os"
	"golang.org/x/sys/windows/registry"
)

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