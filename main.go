package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	appWidth := 400
	appHeight := 400

	wails.Run(&options.App{
		Title:         "Forlife Media Player",
		Width:         appWidth,
		Height:        appHeight,
		DisableResize: true,
		AlwaysOnTop:   true,
		Assets:        assets,
		Bind: []interface{}{
			app,
		},
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
			screens, _ := wailsRuntime.ScreenGetAll(ctx)
			if len(screens) > 0 {
				screen := screens[0]
				x := screen.Size.Width - appWidth - 10
				y := screen.Size.Height - appHeight - 60
				wailsRuntime.WindowSetPosition(ctx, x, y)
			}
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	})
}