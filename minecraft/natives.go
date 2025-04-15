package minecraft

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getNatives(lib ClientJsonLibrary) string {
	archType := "64"
	if strings.Contains(runtime.GOARCH, "386") {
		archType = "32"
	}

	if lib.Natives != nil {
		switch runtime.GOOS {
		case "windows":
			if lib.Natives.Windows != nil {
				return strings.ReplaceAll(*lib.Natives.Windows, "${arch}", archType)
			}
		case "darwin":
		if lib.Natives.Osx != nil {
			return strings.ReplaceAll(*lib.Natives.Osx, "${arch}", archType)
		}
		case "linux":
			if lib.Natives.Linux != nil {
				return strings.ReplaceAll(*lib.Natives.Linux, "${arch}", archType)
			}
		}
	}

	return ""
}

func extractNativesFile(filename, extractPath string, exclude []string) error {
	err := os.MkdirAll(extractPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create extract path: %w", err)
	}

	r, err := zip.OpenReader(filename)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		skip := false
		for _, prefix := range exclude {
			if strings.HasPrefix(f.Name, prefix) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		fullPath := filepath.Join(extractPath, f.Name)

		if !strings.HasPrefix(fullPath, filepath.Clean(extractPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fullPath)
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fullPath, f.Mode())
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
	}

	return nil
}
