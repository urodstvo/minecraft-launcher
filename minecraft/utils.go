package minecraft

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/ulikunitz/xz/lzma"
)

var (
	_versionCache string
	_versionOnce  sync.Once
)

func getUserAgent() string {
	return "MinecraftDownloader/1.0"
}

func getSHA1Hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func checkPathInsideMinecraftDirectory(minecraftDir, path string) error {
	absMinecraftDir, err := filepath.Abs(minecraftDir)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(absPath, absMinecraftDir) {
		return errors.New("path is outside the Minecraft directory")
	}
	return nil
}

func downloadFile(url, path, minecraftDir string, sha1Hash string, overwrite, compressed bool) error {
	if minecraftDir != "" {
		err := checkPathInsideMinecraftDirectory(minecraftDir, path)
		if err != nil {
			return err
		}
	}
	
	if _, err := os.Stat(path); err == nil && !overwrite {
		if sha1Hash == "" {
			return nil
		}

		computedSHA1, err := getSHA1Hash(path)
		if err != nil {
			return err
		}

		if computedSHA1 == sha1Hash {
			return nil
		}
	}

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}


	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()


	if compressed {
		reader, err := lzma.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("error creating lzma reader: %v", err)
		}
		_, err = io.Copy(file, reader)
		if err != nil {
			return err
		}
	} else {
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}
	

	if sha1Hash != "" {
		computedSHA1, err := getSHA1Hash(path)
		if err != nil {
			return err
		}

		if computedSHA1 != sha1Hash {
			return fmt.Errorf("invalid checksum: expected %s, got %s", sha1Hash, computedSHA1)
		}
	}

	return nil
}

// getOSVersion tried to return System.getProperty("os.version") from Java.
func getOSVersion() string {
	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("cmd", "/C", "ver").Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	case "darwin":
		return ""
	default:
		out, err := exec.Command("uname", "-r").Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}
}

func parseSingleRule(rule ClientJsonRule, options *MinecraftOptions) bool {
	var returnValue bool
	if rule.Action == "allow" {
		returnValue = false
	} else if rule.Action == "disallow" {
		returnValue = true
	}


	switch rule.Os.Name {
		case "windows":
			if runtime.GOOS != "windows" {
				return returnValue
			}
		case "osx":
			if runtime.GOOS != "darwin" {
				return returnValue
			}
		case "linux":
			if runtime.GOOS != "linux" {
				return returnValue
			}
		}

	if rule.Os.Arch == "x86" && runtime.GOARCH != "386" {
		return returnValue
	}		
	if matched, _ := regexp.MatchString(rule.Os.Version, getOSVersion()); !matched {
		return returnValue
	}
	
	if rule.Features != nil {
		if !options.CustomResolution && rule.Features.HasCustomResolution  {
			return returnValue
		}
		if !options.Demo && rule.Features.IsDemoUser {
			return returnValue
		}
		if options.QuickPlayPath == nil && rule.Features.HasQuickPlaysSupport {
			return returnValue
		}
		if options.QuickPlaySingleplayer == nil && rule.Features.IsQuickPlaySingleplayer {
			return returnValue
		}
		if options.QuickPlayMultiplayer == nil && rule.Features.IsQuickPlayMultiplayer {
			return returnValue
		}
		if options.QuickPlayRealms == nil && rule.Features.IsQuickPlayRealms {
			return returnValue
		}
		
	}
	
	return !returnValue
}

func parseRuleList(rules []ClientJsonRule, options *MinecraftOptions) bool {
	for _, rule := range rules {
		if !parseSingleRule(rule, options) {
			return false
		}
	}
	return true
}


func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func getClientJson(versionID string, mcDir string) (*ClientJson, error) {
	localPath := filepath.Join(mcDir, "versions", versionID, versionID+".json")

	// Читаем файл
	data, err := readJSON[ClientJson](localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client json: %w", err)
	}

	// Если есть наследование, подгружаем родительский JSON
	if data.InheritsFrom != "" {
		parentData, err := getClientJson(data.InheritsFrom, mcDir)
		if err != nil {
			return nil, fmt.Errorf("failed to inherit from %s: %w", data.InheritsFrom, err)
		}
		// Можешь тут объединить поля из parentData и data вручную
		data = inheritJson(parentData, data)
	}

	return data, nil
}

func readJSON[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data T
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func inheritJson(originalData, parentData *ClientJson) *ClientJson {
	// Создаём карту для проверки существующих библиотек по имени без версии
	originalLibs := make(map[string]bool)
	for _, lib := range originalData.Libraries {
		libName := getLibNameWithoutVersion(lib)
		originalLibs[libName] = true
	}

	// Наследуем библиотеки, которых ещё нет в оригинале
	mergedLibs := originalData.Libraries
	for _, lib := range parentData.Libraries {
		libName := getLibNameWithoutVersion(lib)
		if !originalLibs[libName] {
			mergedLibs = append(mergedLibs, lib)
		}
	}

	// Объединяем структуры (можно углубить при необходимости)
	merged := *parentData
	merged.Libraries = mergedLibs

	// Переписываем остальные поля из originalData (если заданы)
	if originalData.MainClass != "" {
		merged.MainClass = originalData.MainClass
	}
	if originalData.Assets != "" {
		merged.Assets = originalData.Assets
	}
	if originalData.JavaVersion.Component != "" {
		merged.JavaVersion = originalData.JavaVersion
	}
	if originalData.Logging.Client.File.Id != "" {
		merged.Logging = originalData.Logging
	}
	if originalData.Arguments.Game != nil {
		merged.Arguments.Game = append(merged.Arguments.Game, parentData.Arguments.Game...)
	}

	return &merged
}

func getLibNameWithoutVersion(lib ClientJsonLibrary) string {
	parts := strings.Split(lib.Name, ":")
	if len(parts) >= 2 {
		return parts[0] + ":" + parts[1] // groupId:artifactId
	}
	return lib.Name
}

func getClasspathSeparator() string {
	if runtime.GOOS == "windows" {
		return ";"
	}
	return ":"
}

func getLibraryPath(name string, path string) string {
	libPath := filepath.Join(path, "libraries")
	parts := strings.Split(name, ":")

	basePath := parts[0]
	libName := parts[1]
	version := parts[2]

	for _, part := range strings.Split(basePath, ".") {
		libPath = filepath.Join(libPath, part)
	}

	var fileEnd string
	versionParts := strings.Split(version, "@")
	if len(versionParts) == 2 {
		version = versionParts[0]
		fileEnd = versionParts[1]
	} else {
		fileEnd = "jar"
	}

	filenameParts := []string{fmt.Sprintf("%s-%s", libName, version)}
	for _, part := range parts[3:] {
		filenameParts = append(filenameParts, fmt.Sprintf("-%s", part))
	}
	filename := fmt.Sprintf("%s.%s", strings.Join(filenameParts, ""), fileEnd)

	libPath = filepath.Join(libPath, libName, version, filename)
	return libPath
}

func getLibraryVersion() string {
	_versionOnce.Do(func() {		
		filePath := "../.version" 
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading .version:", err)
			_versionCache = "unknown"
			return
		}
		_versionCache = string(data)
	})

	return _versionCache
}