package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urodstvo/minecraft-launcher/minecraft"
)

func EnsureFileExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		file, createErr := os.Create(path)
		if createErr != nil {
			return createErr
		}
		defer file.Close()
		return nil
	}
	return err
}

func LoadCacheToMinecraftOptions(cache launcherCache, mc *minecraft.MinecraftOptions) {
	if cache.Settings != nil {
		settings := cache.Settings
		mc.GameDirectory = settings.GameDirectory
		mc.CustomResolution = settings.ResolutionWidth > 0 && settings.ResolutionHeight > 0

		if settings.ResolutionWidth > 0 {
			mc.ResolutionWidth = strconv.Itoa(settings.ResolutionWidth)
		}

		if settings.ResolutionHeight > 0 {
			mc.ResolutionHeight = strconv.Itoa(settings.ResolutionHeight)
		}

		var jvmArgs []string

		if settings.AllocatedRAM > 0 {
			ramStr := fmt.Sprintf("-Xmx%vM", settings.AllocatedRAM)
			jvmArgs = append(jvmArgs, ramStr)
		}

		if settings.JVMArguments != "" {
			jvmArgs = append(jvmArgs, strings.Fields(settings.JVMArguments)...)
		}

		if len(jvmArgs) > 0 {
			mc.JvmArguments = jvmArgs
		}
	}
}