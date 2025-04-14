package main

import (
	"github.com/urodstvo/minecraft-launcher/minecraft"
)

func main(){
	// logger := utils.Logger{}
	// ui := &launcher.UIView{
	// 	WV: webview.New(false),
	// 	Logger: logger,
	// }
	// l := launcher.Launcher{
	// 	Logger: logger,
	// 	Version: "0.0.0-alpha",
	// 	UI: ui,
	// 	Minecraft: &launcher.Minecraft{},
	// 	Server: &launcher.Server{},
	// 	System: &launcher.System{},
	// }
	// l.Run()
	m := minecraft.Minecraft{
		Config: struct{Directory string}{
			Directory: "Z:\\Games\\Minecraft",
		},
	}
	m.InstallVersion("1.21.5")
}