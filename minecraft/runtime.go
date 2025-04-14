package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
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

func getJVMRuntimes() ([]string, error) {
	resp, err := http.Get(JVM_MANIFEST_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to get JVM manifest: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var manifest map[string]map[string]any
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse JVM manifest: %w", err)
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

func getInstalledJVMRuntimes(minecraftDir string) ([]string, error) {
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


func getJvmRuntimeInformation(jvmVersion string) (*JVMRuntimeInformation, error) {
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

func getVersionRuntimeInformation(versionID string, mcDir string) (*VersionRuntimeInformation, error) {
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

func installRuntimeFile(key string, value _PlatformManifestJsonFile, basePath string, minecraftDirectory string, fileList *[]string, mutex *sync.Mutex) error {
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

		if err := downloadFile(downloadURL, currentPath, "", sha1, false, compressed); err != nil {
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

func (m *Minecraft) InstallJVMRuntime(jvmVersion string) error {
	platform := getJVMPlatform()
	runtimePath := filepath.Join(m.Config.Directory, "runtime", jvmVersion, platform, jvmVersion)

	manifestData, err := fetchManifestData(JVM_MANIFEST_URL)
	if err != nil {
		return fmt.Errorf("failed to get JVM manifest: %w", err)
	}

	runtimeList, ok := manifestData[platform][jvmVersion]
	if !ok || len(runtimeList) == 0 {
		return fmt.Errorf("JVM runtime not found or unsupported for platform: %s", jvmVersion)
	}

	platformManifest, err := fetchPlatformManifest(runtimeList[0].Manifest.Url)
	if err != nil {
		return err
	}


	basePath := path.Join(m.Config.Directory, "runtime", jvmVersion, platform, jvmVersion)

	var fileList []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for path, file := range platformManifest.Files {
		sem <- struct{}{}
		wg.Add(1)
		
		go func(p string, f _PlatformManifestJsonFile) {
			defer wg.Done()
			defer func() { <-sem }() 
			err := installRuntimeFile(p, f, basePath, m.Config.Directory, &fileList, &mu)
			if err != nil {
				fmt.Printf("Error installing file %s: %v\n", p, err)
			}
		}(path, file)
	}

	wg.Wait()

	versionPath := filepath.Join(m.Config.Directory, "runtime", jvmVersion, platform, ".version")
	if err := os.WriteFile(versionPath, []byte(runtimeList[0].Version.Name), 0644); err != nil {
		return err
	}

	sha1Path := filepath.Join(m.Config.Directory, "runtime", jvmVersion, platform, fmt.Sprintf("%s.sha1", jvmVersion))
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

	return nil
}

func fetchManifestData(url string) (RuntimeListJson, error) {
	resp, err := http.Get(url)
	if err != nil {
		return RuntimeListJson{}, fmt.Errorf("failed to fetch manifest data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RuntimeListJson{}, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RuntimeListJson{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var manifestData RuntimeListJson
	if err := json.Unmarshal(body, &manifestData); err != nil {
		return RuntimeListJson{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return manifestData, nil
}

func fetchPlatformManifest(url string) (*PlatformManifestJson, error) {	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch platform manifest data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

		body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	
	var platformManifest PlatformManifestJson
	if err := json.Unmarshal(body, &platformManifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &platformManifest, nil
}