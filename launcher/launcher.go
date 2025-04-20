package launcher

import (
	"fmt"
	"os/exec"
	"runtime"

	"urodstvo-launcher/minecraft"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Launcher struct {
	M minecraft.MinecraftOptions
	cache *launcherCache

	clientId string
	redirectURI string
}

func NewLauncher() *Launcher {
	mc := minecraft.MinecraftOptions{
		LauncherVersion: minecraft.GetLibraryVersion(),
		LauncherName: "Minecraft Launcher by urodstvo.",
	}

	cache := newCache()
	LoadCacheToMinecraftOptions(*cache, &mc)

	return &Launcher{
		M: mc,
		cache: cache,
		clientId: "",
		redirectURI: "",
	}
}

func (a *Launcher) GetLauncherVersion() string {
	return minecraft.GetLibraryVersion()
}

func (a *Launcher) GetMinecraftVersions() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetVersionList()
}

func (a *Launcher) GetLastPlayedVersion() *minecraft.MinecraftVersionInfo {
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

func (a *Launcher) GetInstalledVersion() ([]minecraft.MinecraftVersionInfo, error) {
	return minecraft.GetInstalledVersions(a.M.GameDirectory)
}

func (a *Launcher) OpenMinecraftDirectory() {
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

func (a *Launcher) GetTotalRAM() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total / 1024 / 1024, nil
}

func (a *Launcher) SaveLauncherSettings(settings LauncherSettings) error  {
	a.cache.Settings = &settings
	err := a.cache.Save()
	if err != nil {
		return err
	}

	LoadCacheToMinecraftOptions(*a.cache, &a.M)

	return nil
}

func (a *Launcher) GetLauncherSettings() LauncherSettings {
	return *a.cache.Settings
}

type StartOptions struct {
	Version *minecraft.MinecraftVersionInfo `json:"version"`
}

func (a *Launcher) StartMinecraft(opts StartOptions) {
	a.cache.LastPlayedVersion = opts.Version	
	a.cache.Save()

	fmt.Println(a.M)
	fmt.Println(a.cache)
}

func (a *Launcher) ChooseDirectory() (string, error) {
	dialog := application.OpenFileDialog().CanChooseDirectories(true).CanChooseFiles(false).CanCreateDirectories(true)
	dialog.SetTitle("Select Directory")

	return dialog.PromptForSingleSelection()
}

func (a *Launcher) CreateFreeAccount(username string) {
	newAccount := LauncherAccount{
		Name: username,
		Id: uuid.New().String(),
		Type: "free",
	}
	a.cache.Accounts = append(a.cache.Accounts, newAccount)
	a.cache.Save()

	if len(a.cache.Accounts) == 1 {
		a.SelectAccount(newAccount.Id)
	}
}

// func (a *App) GetSecureMicrosoftLoginURL() (string, error) {
// 	url, state, code, err := minecraft.GetSecureLoginData(a.clientId, a.redirectURI, nil)
// 	return url, err
// }

func (a *Launcher) SelectAccount(id string) {
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

type AccountsInfo struct {
	SelectedAccount string `json:"selectedAccount,omitempty"`
	Accounts []LauncherAccount `json:"accounts,omitempty"`
}

func (a *Launcher) GetAccounts() AccountsInfo{
	return AccountsInfo{
		SelectedAccount: a.cache.SelectedAccount,
		Accounts: a.cache.Accounts,
	}
}