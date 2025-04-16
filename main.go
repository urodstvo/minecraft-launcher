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
	m := minecraft.GenerateTestOptions()
	m.GameDirectory = "Z:\\Games\\Minecraft"
	callback := &minecraft.Callback{
		Progress: func(message string) {
			fmt.Printf("[Progress] - %s\n", message)
		},
		Status: func(message string) {
			fmt.Printf("[Status] - %s\n", message)
		},
		Max: func(message string) {
			fmt.Printf("[Max] - %s\n", message)
		},
	}
	err := minecraft.InstallMinecraftVersion("1.21.3", m, callback)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
		return
	} else {
		fmt.Println("Minecraft version in installed")
	}
	command, err := minecraft.GetMinecraftCommand("1.21.3", m)
	if err != nil {
		fmt.Println(err)
		return
	} 
	// fmt.Println(command)	
	javaPath := command[0]
	javaArgs := command[1:]
	cmd := exec.Command(javaPath, javaArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Println("Запуск Minecraft...")

	// Запускаем
	err = cmd.Run()
	if err != nil {
		fmt.Println("Ошибка при запуске:", err)
	}
}