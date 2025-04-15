package minecraft

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// download library for minecraft
func (m *minecraftConfig) downloadLibrary(id string, lib ClientJsonLibrary) error {
	if len(lib.Rules) > 0 && !parseRuleList(lib.Rules, nil) {
		return nil
	}

	currentPath := filepath.Join(m.Config.Directory, "libraries")
	downloadURL := "https://libraries.minecraft.net"

	if lib.Url != nil {
		downloadURL = *lib.Url
	}

	parts := strings.Split(lib.Name, ":")
	// if len(parts) != 3 {
	// 	return fmt.Errorf("invalid library name format: %s", lib.Name)
	// }

	libPath := parts[0]
	name := parts[1]
	version := parts[2]

	for _, part := range strings.Split(libPath, ".") {
		currentPath = filepath.Join(currentPath, part)
		downloadURL = fmt.Sprintf("%s/%s", downloadURL, part)
	}

	fileEnd := "jar"
	versionParts := strings.Split(version, "@")
	if len(versionParts) == 2 {
		version = versionParts[0]
		fileEnd = versionParts[1]
	}

	jarFileName := fmt.Sprintf("%s-%s.%s", name, version, fileEnd)
	downloadURL = fmt.Sprintf("%s/%s/%s", downloadURL, name, version)
	currentPath = filepath.Join(currentPath, name, version)

	native := getNatives(lib)

	if native != "" {
		jarFileName = fmt.Sprintf("%s-%s-%s.jar", name, version, native)
	}

	err := downloadFile(downloadURL+"/"+jarFileName, filepath.Join(currentPath, jarFileName), m.Config.Directory, "", false, false)
	if err != nil {
		return fmt.Errorf("error downloading library %s: %w", lib.Name, err)
	}

	if lib.Extract != nil && len(lib.Extract.Exclude) > 0 {
		extractNativesFile(filepath.Join(currentPath, jarFileName), filepath.Join(m.Config.Directory, "versions", id, "natives"), lib.Extract.Exclude)
	}

	if len(lib.Downloads.Artifact.Url) > 0 && len(lib.Downloads.Artifact.Path) > 0 {
		err = downloadFile(lib.Downloads.Artifact.Url, filepath.Join(m.Config.Directory, "libraries", lib.Downloads.Artifact.Path), m.Config.Directory, "", false, false)
		if err != nil {
			return fmt.Errorf("error downloading artifact for library %s: %w", lib.Name, err)
		}
	}

	if native != "" && len(lib.Downloads.Classifiers[native].Url) > 0 {
		err = downloadFile(lib.Downloads.Classifiers[native].Url, filepath.Join(currentPath, jarFileName), "", "", false, false)
		if err != nil {
			return fmt.Errorf("error downloading native classifier for library %s: %w", lib.Name, err)
		}
		extractNativesFile(filepath.Join(currentPath, jarFileName), filepath.Join(m.Config.Directory, "versions", id, "natives"), lib.Extract.Exclude)
	}

	return nil
}

// install all necessary libs
func (m *minecraftConfig) installLibraries(id string, libraries []ClientJsonLibrary) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for _, lib := range libraries {
		wg.Add(1)
		sem <- struct{}{} 
		go func(lib ClientJsonLibrary) {
			defer wg.Done()
			if err := m.downloadLibrary(id, lib); err != nil {
				fmt.Printf("Error downloading library %s: %v\n", lib.Name, err)
			}
			<-sem
		}(lib)
	}

	wg.Wait()

	return nil
}

// download asset by its hash
func (m *minecraftConfig) downloadAsset(filehash string) error {
	url := "https://resources.download.minecraft.net/" + filehash[:2] + "/" + filehash
	assetPath := filepath.Join(m.Config.Directory, "assets", "objects", filehash[:2], filehash)
	err := downloadFile(url, assetPath, "", filehash, false, false)
	if err != nil {
		return fmt.Errorf("error downloading asset %s: %v\n", filehash, err)
	}

	return nil
}

// install all assets
func (m *minecraftConfig) installAssets(data ClientJson) error {
	if data.AssetIndex == nil {
		return nil
	}


	assetIndexPath := filepath.Join(m.Config.Directory, "assets", "indexes", data.Assets+".json")
	err := downloadFile(data.AssetIndex.Url, assetIndexPath, "", data.AssetIndex.Sha1, false, false)
	if err != nil {
		return err
	}

	file, err := os.Open(assetIndexPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var assetsData AssetsJson
	err = json.NewDecoder(file).Decode(&assetsData)
	if err != nil {
		return err
	}

	assets := make([]string, 0, len(assetsData.Objects))
	for _, obj := range assetsData.Objects {
		assets = append(assets, obj.Hash)
	}

	c := make(chan uint, maxWorkers) 
	var wg sync.WaitGroup

	for _, filehash := range assets {
		wg.Add(1)

		go func(filehash string) {
			defer wg.Done()

			c <- 1
			defer func() { <-c }()

			err := m.downloadAsset(filehash)
			if err != nil {
				fmt.Println("Error downloading asset:", err)
			}

		}(filehash)
	}

	wg.Wait()

	return nil
}

// install the given version
func (m *minecraftConfig) doVersionInstall(versionID string, url, sha1 string) error {
	versionDir := filepath.Join(m.Config.Directory, "versions", versionID)
	versionJsonPath := filepath.Join(versionDir, versionID+".json")

	if url != "" {
		if err := os.MkdirAll(versionDir, 0755); err != nil {
			return fmt.Errorf("error while creating version: %w", err)
		}
		if err := downloadFile(url, versionJsonPath, m.Config.Directory, sha1, false, false); err != nil {
			return fmt.Errorf("download error of version.json: %w", err)
		}
	}

	var versionData ClientJson
	file, err := os.Open(versionJsonPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	if err := json.Unmarshal(fileContent,  &versionData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	

	// if versionData.InheritsFrom != "" {
	// 	if err := m.InstallMinecraftVersion(versionData.InheritsFrom); err != nil {}
	// 	versionData = inheritJSON(versionData, m.Config.Directory)
	// }


	if err := m.installLibraries(versionData.Id, versionData.Libraries); err != nil {
		return fmt.Errorf("error while installing libraries: %w", err)
	}

	if err := m.installAssets(versionData); err != nil {
		return fmt.Errorf("error while installing assets: %w", err)
	}

	if versionData.Logging.Client.File.Url != "" {
		logFilePath := filepath.Join(m.Config.Directory, "assets", "log_configs", versionData.Logging.Client.File.Id)
		if err := downloadFile(versionData.Logging.Client.File.Url, logFilePath, "", versionData.Logging.Client.File.Sha1, false, false); err != nil {
			return fmt.Errorf("error download log config: %w", err)
		}
	}

	if versionData.Downloads.Client.Url != "" {
		jarPath := filepath.Join(versionDir, versionData.Id+".jar")
		if err := downloadFile(versionData.Downloads.Client.Url, jarPath, "", versionData.Downloads.Client.Sha1, false, false); err != nil {
			return fmt.Errorf("error download client jar: %w", err)
		}
	}

	// jarPath := filepath.Join(versionDir, versionData.Id+".jar")
	// if _, err := os.Stat(jarPath); os.IsNotExist(err) && versionData.InheritsFrom != "" {
	// 	inheritJarPath := filepath.Join(m.Config.Directory, "versions", versionData.InheritsFrom, versionData.InheritsFrom+".jar")
	// 	if err := checkPathInsideMinecraftDirectory(m.Config.Directory, inheritJarPath); err != nil {
	// 		return err
	// 	}
	// 	if err := copyFile(inheritJarPath, jarPath); err != nil {
	// 		return fmt.Errorf("Error copy from parent jar: %w", err)
	// 	}
	// }

	
	if versionData.JavaVersion.Component != "" {
		if err := m.installJVMRuntime(versionData.JavaVersion.Component); err != nil {
			return fmt.Errorf("не удалось установить Java Runtime: %w", err)
		}
	}

	fmt.Println("Version installed:", versionData.Id)
	return nil
}


func (m *minecraftConfig) InstallMinecraftVersion(versionId string) error {
	// versionJsonPath := filepath.Join(m.Config.Directory, "versions", versionId, versionId+".json")
	// if _, err := os.Stat(versionJsonPath); err == nil {
	// 	fmt.Println("Version already installed:", versionId)
	// 	return nil
	// }

	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest_v2.json")
	if err != nil {
		return fmt.Errorf("failed to fetch version list: %w", err)
	}
	defer resp.Body.Close()

	var versionList VersionListManifestJson
	if err := json.NewDecoder(resp.Body).Decode(&versionList); err != nil {
		return fmt.Errorf("failed to decode version list: %w", err)
	}

	for _, version := range versionList.Versions {
		if version.Id == versionId {
			err := m.doVersionInstall(versionId, version.Url, "")
			if err != nil {
				return fmt.Errorf("failed to install version %s: %w", versionId, err)
			}
			return nil
		}
	}

	return fmt.Errorf("version not found")
}