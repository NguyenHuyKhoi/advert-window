package main

import (
	"context"
)

const (
	AppVersionInt    = 2
	BASE_URL         = "https://package-shepherd-knew-proposals.trycloudflare.com/api"
	CHECK_UPDATE_URL = BASE_URL + "/app-settings/laptop"
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