package main

import (
	"fmt"
	"os"
	"os/exec"

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
	command, _ := m.GetMinecraftCommand("1.21.5", m.GenerateTestOptions())
	// fmt.Println(cmd)
	javaPath := command[0]
	javaArgs := command[1:]
	cmd := exec.Command(javaPath, javaArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Println("Запуск Minecraft...")

	// Запускаем
	err := cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при запуске:", err)
	}
}