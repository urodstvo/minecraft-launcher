package main

import (
	"fmt"

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
	m := minecraft.NewAPI(minecraft.Opts{
		MinecraftDirectory: "Z:\\Games\\Minecraft",
		})
	// m.InstallMinecraftVersion("1.21.5")
	cmd, _ := m.GetMinecraftCommand("1.21.5", m.GenerateTestOptions())
	fmt.Println(cmd)
}