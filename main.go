package main

import (
	"embed"
	"fmt"
	"log"

	"urodstvo-launcher/auth"
	"urodstvo-launcher/launcher"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {	
	authService := auth.NewAuthService()
	launcherService := launcher.NewLauncherService()

	app := application.New(application.Options{
		Name: "Minecraft Launcher",
		Description: "by urodstvo",
		Services: []application.Service{
			application.NewService(launcherService, application.ServiceOptions{}),
			application.NewService(authService, application.ServiceOptions{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		ErrorHandler: func(err error) {
			fmt.Println(err.Error())
		},
		WarningHandler: func(s string) {
			fmt.Println(s)
		},
		PanicHandler: func(a any) {
			fmt.Printf("%v\n", a)
		},
	})

	launcherService.SetApp(app)
	authService.SetApp(app)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}