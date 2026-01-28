package main

import (
	"context"
)

const (
	AppVersionInt           = 3
	BASE_URL                = "https://almost-tend-excitement-kijiji.trycloudflare.com/api"
	WINDOW_CHECK_UPDATE_URL = BASE_URL + "/app-settings/window"
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
	go a.silentUpdate()
}

func (a *App) GetVersion() int {
	return AppVersionInt
}

func (a *App) GetDeviceInfo() DeviceInfo {
	return getDeviceInfo()
}
