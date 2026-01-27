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
    // Hàm này sẽ tự động gọi logic từ app_windows.go (nếu là Windows)
    // Hoặc từ app_others.go (nếu là Mac/Linux)
    a.enableAutoStart()
}

func (a *App) GetDeviceInfo() DeviceInfo {
    return getDeviceInfo()
}