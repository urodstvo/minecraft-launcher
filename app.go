package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/urodstvo/minecraft-launcher/minecraft"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	mc minecraft.MinecraftOptions
	cache *launcherCache
}

func NewApp() *App {
	mc := minecraft.MinecraftOptions{
		LauncherVersion: minecraft.GetLibraryVersion(),
		LauncherName: "Minecraft Launcher by urodstvo.",
	}

	cache := newCache()
	LoadCacheToMinecraftOptions(*cache, &mc)

	return &App{
		mc: mc,
		cache: cache,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetLauncherVersion() string {
	return minecraft.GetLibraryVersion()
}

func (a *App) GetMinecraftVersions() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetVersionList()
}

func (a *App) GetLastPlayedVersion() *minecraft.MinecraftVersionInfo {
 if a.cache.LastPlayedVersion == nil {
	var found *minecraft.MinecraftVersionInfo
	v, _ := minecraft.GetLatestVersion()
	l, _ := minecraft.GetVersionList()
	for _, version := range l {
		if version.Id == v.Release {
			found = &version
			break
		}
	}
	return found
 }
 return a.cache.LastPlayedVersion
}

func (a *App) GetInstalledVersion() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetInstalledVersions(a.mc.GameDirectory)
}

func (a *App) OpenMinecraftDirectory() {
	dir := minecraft.GetMinecraftDirectory()
	if a.mc.GameDirectory != "" {
		dir = a.mc.GameDirectory
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	case "linux":
		cmd = exec.Command("xdg-open", dir)
	default:
		return
	}

	_ = cmd.Start()
}

func (a *App) GetTotalRAM() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total / 1024 / 1024, nil
}

func (a *App) SaveLauncherSettings(settings LauncherSettings) error  {
	a.cache.Settings = &settings
	err := a.cache.Save()
	if err != nil {
		return err
	}

	LoadCacheToMinecraftOptions(*a.cache, &a.mc)

	return nil
}

func (a *App) GetLauncherSettings() LauncherSettings {
	return *a.cache.Settings
}

type StartOptions struct {
	Version *minecraft.MinecraftVersionInfo `json:"version"`
}

func (a *App) StartMinecraft(opts StartOptions) {
	a.cache.LastPlayedVersion = opts.Version	
	a.cache.Save()

	fmt.Println(a.mc)
	fmt.Println(a.cache)
}

func (a *App) ChooseDirectory() (string, error) {
	options := wruntime.OpenDialogOptions{
		Title: "Choose Directory",
	}
	dir, err := wruntime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		return "", err
	}
	return dir, nil
}