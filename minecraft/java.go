package minecraft

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func findJavaDirectory(path string) ([]string, error) {
	var javaList []string

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return javaList, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		currentEntry := filepath.Join(path, entry.Name())
		javaBinary := filepath.Join(currentEntry, "bin", "java")
		if runtime.GOOS == "windows" {
			javaBinary += ".exe"
		}

		if _, err := os.Stat(javaBinary); err == nil {
			javaList = append(javaList, currentEntry)
		}
	}

	return javaList, nil
}

func findSystemJavaVersions(additionalDirectories []string) ([]string, error) {
	var javaList []string
	var dirsToSearch []string

	switch runtime.GOOS {
	case "windows":
		dirsToSearch = append(dirsToSearch,
			`C:\Program Files (x86)\Java`,
			`C:\Program Files\Java`,
		)
	case "linux":
		dirsToSearch = append(dirsToSearch,
			"/usr/lib/jvm",
			"/usr/lib/sdk",
		)
	// Note: macOS intentionally left out (unsupported like original Python version)
	}

	dirsToSearch = append(dirsToSearch, additionalDirectories...)

	for _, dir := range dirsToSearch {
		found, err := findJavaDirectory(dir)
		if err != nil {
			// Optional: log or skip on error
			continue
		}
		javaList = append(javaList, found...)
	}

	return javaList, nil
}

func getJavaInformation(path string) (JavaInformation, error) {
	binDir := filepath.Join(path, "bin")
	javaExecutable := "java"
	if runtime.GOOS == "windows" {
		javaExecutable = "java.exe"
	}

	javaPath := filepath.Join(binDir, javaExecutable)

	if _, err := os.Stat(javaPath); os.IsNotExist(err) {
		return JavaInformation{}, errors.New(javaPath + " was not found")
	}

	cmd := exec.Command(javaPath, "-showversion")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return JavaInformation{}, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return JavaInformation{}, errors.New("unexpected java output")
	}

	versionRegex := regexp.MustCompile(`version\s+"([\d._]+)"`)
	versionMatch := versionRegex.FindStringSubmatch(lines[0])
	if versionMatch == nil {
		return JavaInformation{}, errors.New("failed to parse java version")
	}

	info := JavaInformation{
		Path:    path,
		Name:    filepath.Base(path),
		Version: versionMatch[1],
		Openjdk: strings.HasPrefix(strings.ToLower(lines[0]), "openjdk"),
	}

	if len(lines) > 2 && strings.Contains(lines[2], "64-Bit") {
		info.Is64bit = true
	}

	info.JavaPath = javaPath
	if runtime.GOOS == "windows" {
		javaw := filepath.Join(binDir, "javaw.exe")
		info.JavawPath = &javaw
	}

	return info, nil
}

func getSystemJavaVersionInformation(additionalDirectories []string) ([]JavaInformation, error) {
	javaPaths, err := findSystemJavaVersions(additionalDirectories)
	if err != nil {
		return nil, err
	}

	var infos []JavaInformation
	for _, path := range javaPaths {
		info, err := getJavaInformation(path)
		if err != nil {
			log.Printf("warning: failed to get info for %s: %v", path, err)
			continue
		}
		infos = append(infos, info)
	}

	return infos, nil
}