package main

import (
	"encoding/json"
	"os"

	"github.com/urodstvo/minecraft-launcher/minecraft"
)

var _launcherCachePath = "launcherCache.json"

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