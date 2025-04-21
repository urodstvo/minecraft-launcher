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

	window *application.WebviewWindow
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

	l.window = window

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

func (l *LauncherService) GetLauncherVersion() string {
	return minecraft.GetLibraryVersion()
}

func (l *LauncherService) GetMinecraftVersions() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetVersionList()
}

func (l *LauncherService) GetLastPlayedVersion() *minecraft.MinecraftVersionInfo {
 if l.cache.LastPlayedVersion == nil {
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
 return l.cache.LastPlayedVersion
}

func (l *LauncherService) GetInstalledVersion() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetInstalledVersions(l.M.GameDirectory)
}

func (l *LauncherService) OpenMinecraftDirectory() {
	dir := minecraft.GetMinecraftDirectory()
	if l.M.GameDirectory != "" {
		dir = l.M.GameDirectory
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

func (l *LauncherService) GetTotalRAM() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total / 1024 / 1024, nil
}

func (l *LauncherService) SaveLauncherSettings(settings LauncherSettings) error  {
	l.cache.Settings = &settings
	err := l.cache.Save()
	if err != nil {
		return err
	}

	LoadCacheToMinecraftOptions(*l.cache, &l.M)

	return nil
}

func (l *LauncherService) GetLauncherSettings() LauncherSettings {
	return *l.cache.Settings
}

func (l *LauncherService) StartMinecraft(version minecraft.MinecraftVersionInfo) bool {
	l.cache.LastPlayedVersion = &version	
	l.cache.Save()

	if l.M.Uuid == "" {
		return false
	}

	callback := &minecraft.Callback{
		Progress: func(message string) {
			l.app.EmitEvent("install:progress", message)
			// fmt.Println("[Progress] - ", message)
		},
		Status: func(message string) {
			l.app.EmitEvent("install:status", message)
			// fmt.Println("[Status] - ", message)
		},
		Max: func(message string) {
			l.app.EmitEvent("install:max", message)
			// fmt.Println("[Max] - ", message)
		},
	}

	err := minecraft.InstallMinecraftVersion(version.Id, l.M, callback)
	if err != nil {
		return false
	}

	command, err := minecraft.GetMinecraftCommand(version.Id, l.M)
	if err != nil {
		return false
	}

	go func(){
		cmd := exec.Command(command[0], command[1:]...)
		l.window.Close()
		cmd.Run()
		l.window.Show()
		time.Sleep(50 * time.Millisecond)
		l.window.Focus()
	}()

	return true
}

func (l *LauncherService) ChooseDirectory() (string, error) {
	dialog := application.OpenFileDialog().CanChooseDirectories(true).CanChooseFiles(false).CanCreateDirectories(true)
	dialog.SetTitle("Select Directory")

	return dialog.PromptForSingleSelection()
}

func (l *LauncherService) SelectAccount(id string) {
	var selected LauncherAccount
	for _, v := range l.cache.Accounts {
		if v.Id == id {
			selected = v
			break
		}
	}

	l.cache.SelectedAccount = selected.Id
	l.M.Uuid = selected.Id 
	l.M.Username = selected.Name
	l.M.Token = selected.AccessToken
	
	l.cache.Save()
}

func (l *LauncherService) DeleteAccount(id string) {
	if l.cache == nil {
		return
	}

	newAccounts := make([]LauncherAccount, 0, len(l.cache.Accounts))
	for _, acc := range l.cache.Accounts {
		if acc.Id != id {
			newAccounts = append(newAccounts, acc)
		}
	}
	l.cache.Accounts = newAccounts

	if l.cache.SelectedAccount == id {
		l.cache.SelectedAccount = ""
	}

	l.cache.Save()
}

type AccountsInfo struct {
	SelectedAccount string `json:"selectedAccount,omitempty"`
	Accounts []LauncherAccount `json:"accounts,omitempty"`
}

func (l *LauncherService) GetAccounts() AccountsInfo{
	return AccountsInfo{
		SelectedAccount: l.cache.SelectedAccount,
		Accounts: l.cache.Accounts,
	}
}