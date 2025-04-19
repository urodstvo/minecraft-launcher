package main

import (
	"embed"
	"log"
	"runtime"
	"time"

	"urodstvo-launcher/launcher"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	app := application.New(application.Options{
		// Name: l.M.LauncherName,
		// Description: l.M.LauncherVersion,
		Name: "Minecraft Launcher",
		Description: "by urodstvo",
		Services: []application.Service{
			application.NewService(launcher.NewLauncher(), application.ServiceOptions{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,

		},
	})

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Minecraft Launcher",
		Width:  1024,
		Height: 768,
		DisableResize: true,
	})

	systemTray := app.NewSystemTray()
	// if runtime.GOOS == "windows" {
	// 	systemTray.SetIcon(icons.DefaultWindowsIcon)
	// }
	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
		systemTray.SetLabel("Minecraft Launcher")

	}

	myMenu := app.NewMenu()
	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systemTray.SetMenu(myMenu)
	systemTray.OnClick(func() {
		if !window.IsVisible() {
			window.Show()
	
			go func() {
				time.Sleep(150 * time.Millisecond)
				if window.IsVisible() && !window.IsMinimised() {
					window.Focus()
				}
			}()
			return
		}
	
		if window.IsMinimised() {
			window.Restore()
			go func() {
				time.Sleep(150 * time.Millisecond)
				if window.IsVisible() {
					window.Focus()
				}
			}()
		}
	})

	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		window.Hide()
		event.Cancel()
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}