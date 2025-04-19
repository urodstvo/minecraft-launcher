package main

import (
	"encoding/json"
	"os"

	"github.com/urodstvo/minecraft-launcher/minecraft"
)

var _launcherCachePath = "launcherCache.json"

type LauncherAccount struct {
	
}

type LauncherSettings struct {
	GameDirectory string `json:"gameDirectory"`
	AllocatedRAM int `json:"allocatedRAM,omitempty"`
	JVMArguments string `json:"jvmArguments,omitempty"`
	ShowAlpha bool `json:"showAlpha"`
	ShowBeta bool `json:"showBeta"`
	ShowSnaphots bool `json:"showSnapshots"`
	ShowOldVersions bool `json:"showOldVersions"`
	ShowOnlyInstalled bool `json:"showOnlyInstalled"`
	ResolutionWidth int `json:"resolutionWidth,omitempty"`
	ResolutionHeight int `json:"resolutionHeight,omitempty"`
}

type launcherCache struct {
	LastPlayedVersion *minecraft.MinecraftVersionInfo `json:"last_played_version"`
	Settings *LauncherSettings `json:"settings"`
}

func newCache() *launcherCache {
	EnsureFileExists(_launcherCachePath)
	l := &launcherCache{}
	l.Load()

	if l.Settings == nil {
		l.Settings = &LauncherSettings{
			GameDirectory: minecraft.GetMinecraftDirectory(),
			AllocatedRAM: 2048,
			JVMArguments: "",
			ShowAlpha: false,
			ShowBeta: false,
			ShowSnaphots: true,
			ShowOldVersions: false,
			ShowOnlyInstalled: false,
		}
		l.Save()
	}
	
	return l
}

func (c *launcherCache) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(_launcherCachePath, data, 0644)
}

func (c *launcherCache) Load() error {
	data, err := os.ReadFile(_launcherCachePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}