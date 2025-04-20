package launcher

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"urodstvo-launcher/minecraft"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type LauncherService struct {
	M minecraft.MinecraftOptions
	cache *launcherCache

	app *application.App
}

func NewLauncherService() *LauncherService {
	mc := minecraft.MinecraftOptions{
		LauncherVersion: minecraft.GetLibraryVersion(),
		LauncherName: "Minecraft Launcher by urodstvo.",
	}

	cache := newCache()
	LoadCacheToMinecraftOptions(*cache, &mc)

	return &LauncherService{
		M: mc,
		cache: cache,
	}
}

func (l *LauncherService) SetApp(app *application.App) {
	l.app = app

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Minecraft Launcher",
		Name: "Minecraft Launcher",
		Width:  1024,
		Height: 768,
		DisableResize: true,
		FullscreenButtonEnabled: false,
	})

	systemTray := app.NewSystemTray()
	systemTray.SetLabel("Minecraft Launcher")

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

	app.OnEvent("auth:microsoft:failed", func(e *application.CustomEvent) {
		fmt.Println("Auth failed:", e.ToJSON())
	})
	
	app.OnEvent("auth:microsoft:success", func(e *application.CustomEvent) {
		fmt.Println("Auth success:", e.Data)
	})

	app.OnEvent("auth:free:success", func(e *application.CustomEvent) {
		username := e.Data.([]any)[0].(string)
		acc := LauncherAccount{
			Name: username,
			Id: uuid.New().String(),
			Type: "free",
		}

		l.cache.Accounts = append(l.cache.Accounts, acc)
		l.SelectAccount(acc.Id)
		l.cache.Save()
	})
}

func (a *LauncherService) GetLauncherVersion() string {
	return minecraft.GetLibraryVersion()
}

func (a *LauncherService) GetMinecraftVersions() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetVersionList()
}

func (a *LauncherService) GetLastPlayedVersion() *minecraft.MinecraftVersionInfo {
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

func (a *LauncherService) GetInstalledVersion() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetInstalledVersions(a.M.GameDirectory)
}

func (a *LauncherService) OpenMinecraftDirectory() {
	dir := minecraft.GetMinecraftDirectory()
	if a.M.GameDirectory != "" {
		dir = a.M.GameDirectory
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

func (a *LauncherService) GetTotalRAM() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total / 1024 / 1024, nil
}

func (a *LauncherService) SaveLauncherSettings(settings LauncherSettings) error  {
	a.cache.Settings = &settings
	err := a.cache.Save()
	if err != nil {
		return err
	}

	LoadCacheToMinecraftOptions(*a.cache, &a.M)

	return nil
}

func (a *LauncherService) GetLauncherSettings() LauncherSettings {
	return *a.cache.Settings
}

type StartOptions struct {
	Version *minecraft.MinecraftVersionInfo `json:"version"`
}

func (a *LauncherService) StartMinecraft(opts StartOptions) {
	a.cache.LastPlayedVersion = opts.Version	
	a.cache.Save()

	fmt.Println(a.M)
	fmt.Println(a.cache)
}

func (a *LauncherService) ChooseDirectory() (string, error) {
	dialog := application.OpenFileDialog().CanChooseDirectories(true).CanChooseFiles(false).CanCreateDirectories(true)
	dialog.SetTitle("Select Directory")

	return dialog.PromptForSingleSelection()
}

func (a *LauncherService) SelectAccount(id string) {
	var selected LauncherAccount
	for _, v := range a.cache.Accounts {
		if v.Id == id {
			selected = v
			break
		}
	}

	a.cache.SelectedAccount = selected.Id
	a.M.Uuid = selected.Id 
	a.M.Username = selected.Name
	a.M.Token = selected.AccessToken
	
	a.cache.Save()
}

func (a *LauncherService) DeleteAccount(id string) {
	if a.cache == nil {
		return
	}

	newAccounts := make([]LauncherAccount, 0, len(a.cache.Accounts))
	for _, acc := range a.cache.Accounts {
		if acc.Id != id {
			newAccounts = append(newAccounts, acc)
		}
	}
	a.cache.Accounts = newAccounts

	if a.cache.SelectedAccount == id {
		a.cache.SelectedAccount = ""
	}

	a.cache.Save()
}

type AccountsInfo struct {
	SelectedAccount string `json:"selectedAccount,omitempty"`
	Accounts []LauncherAccount `json:"accounts,omitempty"`
}

func (a *LauncherService) GetAccounts() AccountsInfo{
	return AccountsInfo{
		SelectedAccount: a.cache.SelectedAccount,
		Accounts: a.cache.Accounts,
	}
}