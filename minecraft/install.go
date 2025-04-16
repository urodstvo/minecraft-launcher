package minecraft

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func downloadLibrary(id string, lib ClientJsonLibrary, mcDir string) error {
	if len(lib.Rules) > 0 && !parseRuleList(lib.Rules, nil) {
		return nil
	}

	currentPath := filepath.Join(mcDir, "libraries")
	libPath := currentPath
	downloadURL := "https://libraries.minecraft.net"

	if lib.Downloads.Artifact != nil {
		downloadURL = lib.Downloads.Artifact.Url
		libPath = filepath.Join(currentPath, lib.Downloads.Artifact.Path)
	}

	native := getNatives(lib)

	nativeDownloadURL := ""
	nativeLibPath := currentPath
	if native != "" && lib.Downloads.Classifiers != nil {
		nativeDownloadURL = lib.Downloads.Classifiers[native].Url
		nativeLibPath = filepath.Join(currentPath, lib.Downloads.Classifiers[native].Path)
	}

	err := downloadFile(downloadURL, libPath, mcDir, "", false)
	if err != nil {
		return fmt.Errorf("error downloading library %s: %w", lib.Name, err)
	}

	if native != "" {
		err := downloadFile(nativeDownloadURL, nativeLibPath, mcDir, "", false)
		if err != nil {
			return fmt.Errorf("error downloading library %s: %w", lib.Name, err)
		}
		extractNativesFile(libPath, filepath.Join(mcDir, "versions", id, "natives"), lib.Extract.Exclude)
	}

	return nil
}

func installLibraries(id string, libraries []ClientJsonLibrary, mcDir string, callback Callback) error {
	var wg sync.WaitGroup
	var progressWG sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)
	progressCh := make(chan int, len(libraries))

	callback.Status("Downloading Libraries...")
	callback.Progress("0")
    callback.Max(strconv.Itoa(len(libraries)))

	progressWG.Add(1)
	go func() {
		defer progressWG.Done()
		completed := 0
		for progress := range progressCh {
			completed += progress
			callback.Progress(strconv.Itoa(completed))
		}
	}()

	for _, lib := range libraries {
		wg.Add(1)
		sem <- struct{}{} 
		go func(lib ClientJsonLibrary) {
			defer wg.Done()
			defer func() { <-sem }()
			downloadLibrary(id, lib, mcDir)
			progressCh <- 1
		}(lib)
	}

	wg.Wait() 
	close(progressCh)
	progressWG.Wait()

	callback.Status("Libraries download complete.")

	return nil
}

func downloadAsset(filehash string, mcDir string) error {
	url := "https://resources.download.minecraft.net/" + filehash[:2] + "/" + filehash
	assetPath := filepath.Join(mcDir, "assets", "objects", filehash[:2], filehash)
	err := downloadFile(url, assetPath, "", filehash, false)
	if err != nil {
		return fmt.Errorf("error downloading asset %s: %v", filehash, err)
	}

	return nil
}

func installAssets(data ClientJson, mcDir string, callback Callback) error {
	if data.AssetIndex == nil {
		return nil
	}

	assetIndexPath := filepath.Join(mcDir, "assets", "indexes", data.Assets+".json")
	err := downloadFile(data.AssetIndex.Url, assetIndexPath, mcDir, data.AssetIndex.Sha1, false)
	if err != nil {
		return err
	}

	assetsData, err := readJSON[AssetsJson](assetIndexPath)
	if err != nil {
		return err
	}

	assets := make([]string, 0, len(assetsData.Objects))
	for _, obj := range assetsData.Objects {
		assets = append(assets, obj.Hash)
	}

	progressCh := make(chan int, len(assets))
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var progressWG sync.WaitGroup

	callback.Status("Downloading Assets...")
	callback.Progress("0")
    callback.Max(strconv.Itoa(len(assets)))
	
	progressWG.Add(1)
	go func() {
		defer progressWG.Done()
		completed := 0
		for progress := range progressCh {
			completed += progress
			callback.Progress(strconv.Itoa(completed))
		}
	}()

	for _, filehash := range assets {
		wg.Add(1)
		sem <- struct{}{} 
		go func(filehash string) {
			defer wg.Done()			
			defer func() { <-sem }()
			downloadAsset(filehash, mcDir)
			progressCh <- 1
		}(filehash)
	}

	wg.Wait()
	close(progressCh)
	progressWG.Wait()

	callback.Status("Assets download complete.")

	return nil
}

func doVersionInstall(versionID string, url, sha1 string, options MinecraftOptions, callback Callback) error {
	mcDir := options.GameDirectory
	versionDir := filepath.Join(mcDir, "versions", versionID)
	versionJsonPath := filepath.Join(versionDir, versionID+".json")

	if url != "" {
		callback.Status("Downloading Version Manifest")
		callback.Progress("0")
		callback.Max("1")

		if err := os.MkdirAll(versionDir, 0755); err != nil {
			return fmt.Errorf("error while creating version: %w", err)
		}
		if err := downloadFile(url, versionJsonPath, mcDir, sha1, false); err != nil {
			return fmt.Errorf("download error of version.json: %w", err)
		}

		callback.Progress("1")
		callback.Status("Version Manifest download complete.")
	}

	versionData, err := readJSON[ClientJson](versionJsonPath)
	if err != nil {
		return err
	}

	if versionData.InheritsFrom != "" {
		InstallMinecraftVersion(versionData.InheritsFrom, options, &callback);
		versionData, _ = inheritJson(versionData, mcDir)
	}

	if err := installLibraries(versionData.Id, versionData.Libraries, mcDir, callback); err != nil {
		return fmt.Errorf("error while installing libraries: %w", err)
	}

	if err := installAssets(versionData, mcDir, callback); err != nil {
		return fmt.Errorf("error while installing assets: %w", err)
	}

	if versionData.Logging.Client.File.Url != "" {
		logFilePath := filepath.Join(mcDir, "assets", "log_configs", versionData.Logging.Client.File.Id)
		if err := downloadFile(versionData.Logging.Client.File.Url, logFilePath, "", versionData.Logging.Client.File.Sha1, false); err != nil {
			return fmt.Errorf("error download log config: %w", err)
		}
	}

	if versionData.Downloads.Client.Url != "" {
		jarPath := filepath.Join(versionDir, versionData.Id+".jar")
		if err := downloadFile(versionData.Downloads.Client.Url, jarPath, "", versionData.Downloads.Client.Sha1, false); err != nil {
			return fmt.Errorf("error download client jar: %w", err)
		}
	}

	jarPath := filepath.Join(versionDir, versionData.Id+".jar")
	if _, err := os.Stat(jarPath); os.IsNotExist(err) && versionData.InheritsFrom != "" {
		inheritJarPath := filepath.Join(mcDir, "versions", versionData.InheritsFrom, versionData.InheritsFrom+".jar")
		if err := checkPathInsideMinecraftDirectory(mcDir, inheritJarPath); err != nil {
			return err
		}
		if err := copyFile(inheritJarPath, jarPath); err != nil {
			return fmt.Errorf("error copy from parent jar: %w", err)
		}
	}
	
	if versionData.JavaVersion.Component != "" {
		if err := installJVMRuntime(versionData.JavaVersion.Component, mcDir, callback); err != nil {
			return fmt.Errorf("error installing Java Runtime: %w", err)
		}
	}
	return nil
}

func InstallMinecraftVersion(versionId string, options MinecraftOptions, callback *Callback) error {
	versionList, err := fetch[VersionListManifestJson]("https://launchermeta.mojang.com/mc/game/version_manifest_v2.json")
	if err != nil {
		return fmt.Errorf("failed to decode version list: %w", err)
	}

	if callback == nil {
		callback = &Callback{
			Progress: func(message string) {},
			Max: func(message string) {},
			Status: func(message string) {},
		}
	}


	for _, version := range versionList.Versions {
		if version.Id == versionId {
			err := doVersionInstall(versionId, version.Url, "", options, *callback)
			if err != nil {
				return fmt.Errorf("failed to install version %s: %w", versionId, err)
			}
			return nil
		}
	}

	return ErrorVersionNotFound
}