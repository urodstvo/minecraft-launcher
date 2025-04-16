package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const JVM_MANIFEST_URL string = "https://launchermeta.mojang.com/v1/products/java-runtime/2ec0cc96c44e5a76b9c8b7c39df7210883d12871/all.json"

func getJVMPlatform() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goos {
	case "windows":
		if goarch == "386" {
			return "windows-x86"
		}
		return "windows-x64"

	case "linux":
		if goarch == "386" {
			return "linux-i386"
		}
		return "linux"

	case "darwin":
		if goarch == "arm64" {
			return "mac-os-arm64"
		}
		return "mac-os"

	default:
		return "gamecore"
	}
}

func GetJVMRuntimes() ([]string, error) {
	manifest, err := fetch[map[string]map[string]any](JVM_MANIFEST_URL)
	if err != nil {		
		return nil, fmt.Errorf("error fetching platform manifest: %v", err)
	}

	platform := getJVMPlatform()
	platformData, ok := manifest[platform]
	if !ok {
		return nil, fmt.Errorf("platform %s not found in manifest", platform)
	}

	runtimes := make([]string, 0, len(platformData))
	for key := range platformData {
		runtimes = append(runtimes, key)
	}

	return runtimes, nil
}

func GetInstalledJVMRuntimes(minecraftDir string) ([]string, error) {
	runtimeDir := filepath.Join(minecraftDir, "runtime")
	entries, err := os.ReadDir(runtimeDir)
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}

	runtimes := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			runtimes = append(runtimes, entry.Name())
		}
	}
	return runtimes, nil
}

func getExecutablePath(jvmVersion, minecraftDir string) string {
	platform := getJVMPlatform()

	basePath := filepath.Join(minecraftDir, "runtime", jvmVersion, platform, jvmVersion)
	javaPath := filepath.Join(basePath, "bin", "java")

	if fileExists(javaPath) {
		return javaPath
	}
	if fileExists(javaPath + ".exe") {
		return javaPath + ".exe"
	}

	javaPath = filepath.Join(basePath, "jre.bundle", "Contents", "Home", "bin", "java")
	if fileExists(javaPath) {
		return javaPath
	}

	return ""
}

func GetJvmRuntimeInformation(jvmVersion string) (*JVMRuntimeInformation, error) {
	platform := getJVMPlatform()

	resp, err := http.Get(JVM_MANIFEST_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	var manifest map[string]map[string][]struct {
		Version struct {
			Name     string `json:"name"`
			Released string `json:"released"`
		} `json:"version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}

	platformData, ok := manifest[platform]
	if !ok {
		return nil, errors.New("JVM version not available for this platform")
	}

	versions, ok := platformData[jvmVersion]
	if !ok {
		return nil, errors.New("JVM version not found")
	}
	if len(versions) == 0 {
		return nil, errors.New("JVM version not available for this platform")
	}

	releasedTime, err := time.Parse(time.RFC3339, versions[0].Version.Released)
	if err != nil {
		return nil, fmt.Errorf("invalid release date format: %w", err)
	}

	return &JVMRuntimeInformation{
		Name:     versions[0].Version.Name,
		Released: releasedTime,
	}, nil
}

func GetVersionRuntimeInformation(versionID string, mcDir string) (*VersionRuntimeInformation, error) {
	data, err := getClientJson(versionID, mcDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load client json: %w", err)
	}

	if data.JavaVersion.Component == "" {
		return nil, nil
	}

	return &VersionRuntimeInformation{
		Name:             data.JavaVersion.Component,
		JavaMajorVersion: data.JavaVersion.MajorVersion,
	}, nil
}

func installRuntimeFile(key string, value platformManifestJsonFile, basePath string, minecraftDirectory string, fileList *[]string, mutex *sync.Mutex) error {
	currentPath := filepath.Join(basePath, key)

	if err := checkPathInsideMinecraftDirectory(minecraftDirectory, currentPath); err != nil {
		return err
	}

	switch value.Type {
	case "file":
		var downloadURL string
		var sha1 string
		var compressed bool

		if lzma, ok := value.Downloads["lzma"]; ok {
			downloadURL = lzma.Url
			sha1 = value.Downloads["raw"].SHA1
			compressed = true
		} else {
			downloadURL = value.Downloads["raw"].Url
			sha1 = value.Downloads["raw"].SHA1
		}

		var err error
		if compressed {
			err = downloadCompressedFile(downloadURL, currentPath, minecraftDirectory, sha1, false); 
		} else {
			err = downloadFile(downloadURL, currentPath, minecraftDirectory, sha1, false); 
		}
		if err != nil {
			return err
		}

		if value.Executable {
			if err := os.Chmod(currentPath, 0755); err != nil {
				return err
			}
		}

		mutex.Lock()
		*fileList = append(*fileList, key)
		mutex.Unlock()

	case "directory":
		if err := os.MkdirAll(currentPath, os.ModePerm); err != nil {
			return err
		}

	case "link":
		targetPath := filepath.Join(basePath, value.Target)
		if err := checkPathInsideMinecraftDirectory(minecraftDirectory, targetPath); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(currentPath), os.ModePerm); err != nil {
			return err
		}

		if err := os.Symlink(value.Target, currentPath); err != nil {
			return err
		}
	}

	return nil
}

func installJVMRuntime(jvmVersion string, mcDir string, callback Callback) error {
	platform := getJVMPlatform()
	runtimePath := filepath.Join(mcDir, "runtime", jvmVersion, platform, jvmVersion)

	manifestData, err := fetch[RuntimeListJson](JVM_MANIFEST_URL)
	if err != nil {		
		return fmt.Errorf("error fetching jvm manifest: %v", err)
	}

	runtimeList, ok := manifestData[platform][jvmVersion]
	if !ok || len(runtimeList) == 0 {
		return fmt.Errorf("JVM runtime not found or unsupported for platform: %s", jvmVersion)
	}

	platformManifest, err := fetch[PlatformManifestJson](runtimeList[0].Manifest.Url)
	if err != nil {		
		return fmt.Errorf("error fetching platform manifest: %v", err)
	}

	basePath := path.Join(mcDir, "runtime", jvmVersion, platform, jvmVersion)

	var fileList []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var progressWG sync.WaitGroup
	sem := make(chan struct{}, maxWorkers) 
	progressCh := make(chan int, len(platformManifest.Files))

	callback.Status("Downloading JVM Runtime Files...")
	callback.Progress("0")
	callback.Max(strconv.Itoa(len(platformManifest.Files)))

	progressWG.Add(1)
	go func() {
		defer progressWG.Done()
		completed := 0
		for progress := range progressCh {
			completed += progress
			callback.Progress(strconv.Itoa(completed))
		}
	}()

	for path, file := range platformManifest.Files {
		wg.Add(1)
		sem <- struct{}{}
		go func(p string, f platformManifestJsonFile) {
			defer wg.Done()
			defer func() { <-sem }()
			installRuntimeFile(p, f, basePath, mcDir, &fileList, &mu)
			progressCh <- 1
		}(path, file)
	}

	wg.Wait() 
	close(progressCh)
	progressWG.Wait()

	callback.Status("JVM Runtime Files download complete.")
	callback.Status("Installing JVM Runtime Files...")
	callback.Progress("0")
    callback.Max("1")

	versionPath := filepath.Join(mcDir, "runtime", jvmVersion, platform, ".version")
	if err := os.WriteFile(versionPath, []byte(runtimeList[0].Version.Name), 0644); err != nil {
		return err
	}

	sha1Path := filepath.Join(mcDir, "runtime", jvmVersion, platform, fmt.Sprintf("%s.sha1", jvmVersion))
	f, err := os.Create(sha1Path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, file := range fileList {
		fullPath := filepath.Join(runtimePath, file)
		stat, _ := os.Stat(fullPath)
		hash, err := getSHA1Hash(fullPath)
		if err != nil {
			return err
		}
		fmt.Fprintf(f, "%s /#// %s %d\n", file, hash, stat.ModTime().UnixNano())
	}

	callback.Progress("1")
	callback.Status("JVM Runtime Files install complete.")

	return nil
}

