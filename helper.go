package main

import "os"

func EnsureFileExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		file, createErr := os.Create(path)
		if createErr != nil {
			return createErr
		}
		defer file.Close()
		return nil
	}
	return err
}