package main

import (
	"context"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.enableAutoStart()
}

// Logic thực tế cho Windows nằm ở app_windows.go
func (a *App) enableAutoStart() {
}

func (a *App) GetDeviceInfo() DeviceInfo {
	return getDeviceInfo()
}