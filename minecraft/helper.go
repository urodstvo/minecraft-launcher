package minecraft

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/ulikunitz/xz/lzma"
)

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

func getOSVersion() string {
	if runtime.GOOS == "windows" {
		maj, min, _ := windows.RtlGetNtVersionNumbers()		
		return fmt.Sprintf("%d.%d", maj, min)
	}

	if runtime.GOOS == "darwin" {		
		return ""
	}

	return runtime.GOARCH 
}

func parseSingleRule(rule ClientJsonRule, options *MinecraftOptions) bool {
	var returnValue bool
	if rule.Action == "allow" {
		returnValue = false
	} else if rule.Action == "disallow" {
		returnValue = true
	}

	if rule.Os.Name != nil {
		switch *rule.Os.Name {
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
	}

	if rule.Os.Arch != nil {
		if *rule.Os.Arch == "x86" && runtime.GOARCH != "386" {
			return returnValue
		}
		if *rule.Os.Arch == "x64" && runtime.GOARCH != "amd64" {
			return returnValue
		}
	}

	if rule.Os.Version != nil {
		if matched, _ := regexp.MatchString(*rule.Os.Version, getOSVersion()); !matched {
			return returnValue
		}
	}


	if rule.Features != nil {
		if !options.CustomResolution && rule.Features.HasCustomResolution != nil && *rule.Features.HasCustomResolution  {
			return returnValue
		}
		if !options.Demo && rule.Features.IsDemoUser != nil && *rule.Features.IsDemoUser {
			return returnValue
		}
		if options.QuickPlayPath == nil && rule.Features.HasQuickPlaysSupport != nil && *rule.Features.HasQuickPlaysSupport {
			return returnValue
		}
		if options.QuickPlaySingleplayer == nil && rule.Features.IsQuickPlaySingleplayer != nil && *rule.Features.IsQuickPlaySingleplayer {
			return returnValue
		}
		if options.QuickPlayMultiplayer == nil && rule.Features.IsQuickPlayMultiplayer != nil && *rule.Features.IsQuickPlayMultiplayer {
			return returnValue
		}
		if options.QuickPlayRealms == nil && rule.Features.IsQuickPlayRealms != nil && *rule.Features.IsQuickPlayRealms {
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

func getClientJson(version string, minecraftDirectory string) (ClientJson, error) {
	localPath := filepath.Join(minecraftDirectory, "versions", version, fmt.Sprintf("%s.json", version))

	if _, err := os.Stat(localPath); err == nil {
		file, err := os.Open(localPath)
		if err != nil {
			return ClientJson{}, err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return ClientJson{}, err
		}

		var clientData ClientJson
		if err := json.Unmarshal(data, &clientData); err != nil {
			return ClientJson{}, err
		}

		if clientData.InheritsFrom != "" {
			clientData, err = inheritJson(clientData, minecraftDirectory)
			if err != nil {
				return ClientJson{}, err
			}
		}

		return clientData, nil
}

versionListURL := "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json"
	resp, err := getRequestsResponseCache(versionListURL)
	if err != nil {
		return ClientJson{}, err
	}
	defer resp.Body.Close()

	var versionList map[string][]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&versionList); err != nil {
		return ClientJson{}, err
	}

	for _, v := range versionList["versions"] {
		if v["id"] == version {
			resp, err := getRequestsResponseCache(v["url"])
			if err != nil {
				return ClientJson{}, err
			}
			defer resp.Body.Close()

			var clientData ClientJson
			if err := json.NewDecoder(resp.Body).Decode(&clientData); err != nil {
				return ClientJson{}, err
			}

			return clientData, nil
		}
	}

	return ClientJson{}, errors.New("version not found: " + version)
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

func inheritJson(originalData ClientJson, path string) (ClientJson, error) {
	inheritVersion := originalData.InheritsFrom


	inheritFilePath := filepath.Join(path, "versions", inheritVersion, inheritVersion+".json")
	inheritFile, err := os.Open(inheritFilePath)
	if err != nil {
		fmt.Println("Failed to open file: %w", err)
		return ClientJson{}, err
	}
	defer inheritFile.Close()

	data, err := io.ReadAll(inheritFile)
	if err != nil {
		return ClientJson{}, err
	}

	var newData ClientJson
	if err := json.Unmarshal(data, &newData); err != nil {
		return ClientJson{}, err
	}

	originalLibs := make(map[string]bool)
	for _, lib := range originalData.Libraries {
		libName := getLibNameWithoutVersion(lib)
		originalLibs[libName] = true
	}

	var libList []ClientJsonLibrary
	for _, lib := range newData.Libraries {
		libName := getLibNameWithoutVersion(lib)
		if _, exists := originalLibs[libName]; !exists {
			libList = append(libList, lib)
		}		
	}
	newData.Libraries = libList

	if originalData.Arguments != nil && newData.Arguments != nil {
		newData.Arguments.Game = append(newData.Arguments.Game, originalData.Arguments.Game...)
		newData.Arguments.Jvm = append(newData.Arguments.Jvm, originalData.Arguments.Jvm...)
	}

	if originalData.Downloads.Client != (clientJsonDownloads{}) {
		newData.Downloads.Client = originalData.Downloads.Client
	}
	if originalData.Downloads.Server != (clientJsonDownloads{}) {
		newData.Downloads.Server = originalData.Downloads.Server
	}

	if originalData.Logging.Client != (clientJsonLogging{}) {
		newData.Logging.Client = originalData.Logging.Client
	}

	if originalData.MainClass != "" {
		newData.MainClass = originalData.MainClass
	}
	if originalData.MinimumLauncherVersion != 0 {
		newData.MinimumLauncherVersion = originalData.MinimumLauncherVersion
	}
	if originalData.ReleaseTime != "" {
		newData.ReleaseTime = originalData.ReleaseTime
	}
	if originalData.Time != "" {
		newData.Time = originalData.Time
	}
	if originalData.Type != "" {
		newData.Type = originalData.Type
	}

	newData.ComplianceLevel = originalData.ComplianceLevel

	return newData, nil
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

	for part := range strings.SplitSeq(basePath, ".") {
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