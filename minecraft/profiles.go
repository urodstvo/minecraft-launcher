package minecraft

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func isVanillaLauncherProfileValid(profile VanillaLauncherProfile) bool {
	if profile.Name == "" {
		return false
	}

	switch profile.VersionType {
	case "latest-release", "latest-snapshot":
	case "custom":
		if profile.Version == nil {
			return false
		}
	default:
		return false
	}

	if profile.GameDirectory != nil && *profile.GameDirectory == "" {
		return false
	}

	if profile.JavaExecutable != nil && *profile.JavaExecutable == "" {
		return false
	}

	for _, arg := range profile.JavaArguments {
		if arg == "" {
			return false
		}
	}

	if res := profile.CustomResolution; res != nil {
		if res.Height <= 0 || res.Width <= 0 {
			return false
		}
	}

	return true
}

func LoadVanillaLauncherProfiles(minecraftDir string) ([]VanillaLauncherProfile, error) {
	filePath := filepath.Join(minecraftDir, "launcher_profiles.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read launcher_profiles.json: %w", err)
	}

	var launcherData VanillaLauncherProfilesJson
	if err := json.Unmarshal(data, &launcherData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var profiles []VanillaLauncherProfile
	for _, value := range launcherData.Profiles {
		var profile VanillaLauncherProfile

		switch value.Type {
		case "latest-release":
			profile.Name = "Latest release"
		case "latest-snapshot":
			profile.Name = "Latest snapshot"
		default:
			profile.Name = value.Name
		}

		switch value.LastVersionID {
		case "latest-release":
			profile.VersionType = "latest-release"
		case "latest-snapshot":
			profile.VersionType = "latest-snapshot"
		default:
			profile.VersionType = "custom"
			profile.Version = &value.LastVersionID
		}

		if value.GameDir != "" {
			profile.GameDirectory = &value.GameDir
		}
		if value.JavaDir != "" {
			profile.JavaExecutable = &value.JavaDir
		}
		if value.JavaArgs != "" {
			profile.JavaArguments = strings.Fields(value.JavaArgs)
		}
		if value.Resolution != nil {
			profile.CustomResolution = value.Resolution
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func VanillaLauncherProfileToMinecraftOptions(profile VanillaLauncherProfile) (MinecraftOptions, error) {
	if !isVanillaLauncherProfileValid(profile) {
		return MinecraftOptions{}, fmt.Errorf("invalid vanilla launcher profile")
	}

	var opts MinecraftOptions

	if profile.GameDirectory != nil {
		opts.GameDirectory = *profile.GameDirectory
	}

	if profile.JavaExecutable != nil {
		opts.ExecutablePath = *profile.JavaExecutable
	}

	if profile.JavaArguments != nil {
		opts.JvmArguments = profile.JavaArguments
	}

	if profile.CustomResolution != nil {
		opts.CustomResolution = true
		opts.ResolutionWidth = strconv.Itoa(profile.CustomResolution.Width)
		opts.ResolutionHeight = strconv.Itoa(profile.CustomResolution.Height)
	}

	return opts, nil
}


func AddVanillaLauncherProfile(minecraftDir string, profile VanillaLauncherProfile) error {
	if !isVanillaLauncherProfileValid(profile) {
		return fmt.Errorf("invalid vanilla launcher profile")
	}

	filePath := filepath.Join(minecraftDir, "launcher_profiles.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read launcher_profiles.json: %w", err)
	}

	var launcher VanillaLauncherProfilesJson
	if err := json.Unmarshal(data, &launcher); err != nil {
		return fmt.Errorf("failed to parse launcher_profiles.json: %w", err)
	}

	var newProfile VanillaLauncherProfilesJsonProfile
	newProfile.Name = profile.Name

	switch profile.VersionType {
	case "latest-release":
		newProfile.LastVersionID = "latest-release"
	case "latest-snapshot":
		newProfile.LastVersionID = "latest-snapshot"
	case "custom":
		newProfile.LastVersionID = derefStr(profile.Version)
	default:
		return fmt.Errorf("unsupported versionType: %s", profile.VersionType)
	}

	if profile.GameDirectory != nil {
		newProfile.GameDir = *profile.GameDirectory
	}
	if profile.JavaExecutable != nil {
		newProfile.JavaDir = *profile.JavaExecutable
	}
	if len(profile.JavaArguments) > 0 {
		newProfile.JavaArgs = strings.Join(profile.JavaArguments, " ")
	}
	if profile.CustomResolution != nil {
		newProfile.Resolution = profile.CustomResolution
	}

	now := time.Now().Format(time.RFC3339)
	newProfile.Created = now
	newProfile.LastUsed = now
	newProfile.Type = "custom"

	var key string
	for {
		key = uuid.NewString()
		if _, exists := launcher.Profiles[key]; !exists {
			break
		}
	}

	launcher.Profiles[key] = newProfile

	out, err := json.MarshalIndent(launcher, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal launcher_profiles.json: %w", err)
	}

	if err := os.WriteFile(filePath, out, 0644); err != nil {
		return fmt.Errorf("failed to write launcher_profiles.json: %w", err)
	}

	return nil
}

func GetVanillaLauncherProfileVersion(profile VanillaLauncherProfile, latestVersionFunc func() map[string]string) (string, error) {
	if !isVanillaLauncherProfileValid(profile) {
		return "", fmt.Errorf("invalid vanilla launcher profile")
	}

	switch profile.VersionType {
	case "latest-release":
		return latestVersionFunc()["release"], nil
	case "latest-snapshot":
		return latestVersionFunc()["snapshot"], nil
	case "custom":
		if profile.Version != nil {
			return *profile.Version, nil
		}
		return "", fmt.Errorf("custom version type but version is nil")
	default:
		return "", fmt.Errorf("unsupported version type: %s", profile.VersionType)
	}
}

func DoVanillaLauncherProfilesExist(minecraftDir string) bool {
	filePath := filepath.Join(minecraftDir, "launcher_profiles.json")
	_, err := os.Stat(filePath)
	return err == nil
}