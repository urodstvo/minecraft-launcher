package minecraft

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	mu sync.Mutex
	_versionCache string
	_versionOnce  sync.Once
	_requestsResponseCache = make(map[string]requestsResponseCache)
	_userAgentCache string
)

func getUserAgent() string {
	mu.Lock()
	defer mu.Unlock()

	if _userAgentCache != "" {
		return _userAgentCache
	}

	versionFilePath := ".version"
	data, err := os.ReadFile(versionFilePath)
	if err != nil {
		return ""
	} 	

	_userAgentCache = "urodstvo-launcher/" + strings.TrimSpace(string(data))
	return _userAgentCache
}

func GetLibraryVersion() string {
	_versionOnce.Do(func() {		
		filePath := ".version" 
		data, err := os.ReadFile(filePath)
		if err != nil {
			_versionCache = "unknown"
			return
		}
		_versionCache = string(data)
	})

	return _versionCache
}

func GetMinecraftDirectory() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, _ := os.UserHomeDir()
			return filepath.Join(homeDir, "AppData", "Roaming", ".minecraft")
		}
		return filepath.Join(appData, ".minecraft")
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Application Support", "minecraft")
	default:
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".minecraft")
	}
}

func getRequestsResponseCache(url string) (*http.Response, error) {
	mu.Lock()
	defer mu.Unlock()

	cache, found := _requestsResponseCache[url]
	if found && time.Since(cache.Datetime).Hours() < 1 {
		return cache.Response, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		_requestsResponseCache[url] = requestsResponseCache{
			Response: resp,
			Datetime: time.Now(),
		}
	}

	return resp, nil
}

func GetLatestVersion() (LatestMinecraftVersions, error) {
	resp, err := getRequestsResponseCache("https://launchermeta.mojang.com/mc/game/version_manifest_v2.json")
	if err != nil {
		return LatestMinecraftVersions{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Latest LatestMinecraftVersions `json:"latest"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return LatestMinecraftVersions{}, err
	}

	return result.Latest, nil
}

func GetVersionList() ([]MinecraftVersionInfo, error) {
	resp, err := getRequestsResponseCache("https://launchermeta.mojang.com/mc/game/version_manifest_v2.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vlist VersionListManifestJson
	if err := json.NewDecoder(resp.Body).Decode(&vlist); err != nil {
		return nil, err
	}

	var res []MinecraftVersionInfo
	for _, v := range vlist.Versions {
		res = append(res, MinecraftVersionInfo{
			Id: v.Id,
			Type: v.Type,
			ReleaseTime: v.ReleaseTime,
			ComplianceLevel: int(v.ComplianceLevel),
		})
	}

	return res, nil
}

func GetInstalledVersions(minecraftDirectory string) ([]MinecraftVersionInfo, error) {
	versionsPath := filepath.Join(minecraftDirectory, "versions")
	dirEntries, err := os.ReadDir(versionsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []MinecraftVersionInfo{}, nil
		}
		return nil, err
	}

	var versionList []MinecraftVersionInfo

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}

		versionID := entry.Name()
		jsonPath := filepath.Join(versionsPath, versionID, versionID+".json")

		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			continue 
		}

		data, err := os.ReadFile(jsonPath)
		if err != nil {
			continue
		}

		var versionData versionListManifestJsonVersion
		if err := json.Unmarshal(data, &versionData); err != nil {
			continue
		}

		releaseTime, err := time.Parse(time.RFC3339, versionData.ReleaseTime)
		if err != nil {
			releaseTime = time.Unix(0, 0)
		}

		versionList = append(versionList, MinecraftVersionInfo{
			Id:              versionData.Id,
			Type:            versionData.Type,
			ReleaseTime:     releaseTime.String(),
			ComplianceLevel: versionData.ComplianceLevel,
		})
	}

	return versionList, nil
}

func GetAvailableVersions(minecraftDirectory string) ([]MinecraftVersionInfo, error) {
	versionList, err := GetVersionList()
	if err != nil {
		return nil, err
	}

	installedVersions, err := GetInstalledVersions(minecraftDirectory)
	if err != nil {
		return nil, err
	}

	installedMap := make(map[string]bool)
	for _, version := range installedVersions {
		installedMap[version.Id] = true
	}

	var combinedVersions []MinecraftVersionInfo

	combinedVersions = append(combinedVersions, versionList...)

	for _, version := range installedVersions {
		if _, exists := installedMap[version.Id]; !exists {
			combinedVersions = append(combinedVersions, version)
		}
	}

	return combinedVersions, nil
}

func GenerateTestOptions() MinecraftOptions {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	
	username := fmt.Sprintf("Player%d", rand.Intn(900)+100) 
	uuidValue := uuid.New().String()                        

	return MinecraftOptions{
		Username: username,
		Uuid:     uuidValue,
		Token:    "", 
		Demo: false,
	}
}

func IsPlatformSupported() bool {
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		return true
	default:
		return false
	}
}

func IsMinecraftInstalled(minecraftDirectory string) bool {
	requiredDirs := []string{"versions", "libraries", "assets"}

	for _, dir := range requiredDirs {
		fullPath := filepath.Join(minecraftDirectory, dir)
		info, err := os.Stat(fullPath)
		if err != nil || !info.IsDir() {
			return false
		}
	}
	return true
}