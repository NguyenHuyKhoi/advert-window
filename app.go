package main

import (
	"context"
)

const AppVersionInt = 1

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

// Thêm lại hàm này để fix lỗi biên dịch TypeScript
func (a *App) GetDeviceInfo() DeviceInfo {
	return getDeviceInfo()
}