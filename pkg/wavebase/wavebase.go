// Copyright 2024, Command Line Inc.
// SPDX-License-Identifier: Apache-2.0

// Package wavebase provides core utilities and base functionality for WaveTerm.
package wavebase

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const WaveTermAppName = "waveterm"
const WaveTermVersion = "v0.1.0"
const WaveTermDirName = ".waveterm"
const WaveTermDevDirName = ".waveterm-dev"

// dirPermissions defines the permission bits used when creating WaveTerm directories.
// Using 0700 instead of 0755 to restrict access to the current user only.
const dirPermissions = 0700

var baseLock sync.Mutex
var waveHomeDir string

// IsDevMode returns true if the application is running in development mode.
func IsDevMode() bool {
	pname := os.Getenv("WAVETERM_DEV")
	return pname != ""
}

// GetWaveHomeDir returns the path to the WaveTerm home directory,
// creating it if it does not exist. The directory is determined once
// and cached for subsequent calls.
func GetWaveHomeDir() (string, error) {
	baseLock.Lock()
	defer baseLock.Unlock()

	if waveHomeDir != "" {
		return waveHomeDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user home directory: %w", err)
	}

	dirName := WaveTermDirName
	if IsDevMode() {
		dirName = WaveTermDevDirName
	}

	waveDir := filepath.Join(homeDir, dirName)
	if err := ensureDir(waveDir); err != nil {
		return "", fmt.Errorf("cannot create wave home directory %q: %w", waveDir, err)
	}

	waveHomeDir = waveDir
	return waveHomeDir, nil
}

// EnsureWaveHomeDir creates the WaveTerm home directory if it does not exist.
func EnsureWaveHomeDir() error {
	_, err := GetWaveHomeDir()
	return err
}

// GetWaveDataDir returns the path to the WaveTerm data directory.
func GetWaveDataDir() (string, error) {
	homeDir, err := GetWaveHomeDir()
	if err != nil {
		return "", err
	}
	dataDir := filepath.Join(homeDir, "data")
	if err := ensureDir(dataDir); err != nil {
		return "", fmt.Errorf("cannot create wave data directory %q: %w", dataDir, err)
	}
	return dataDir, nil
}

// GetWaveLogDir returns the path to the WaveTerm log directory.
func GetWaveLogDir() (string, error) {
	homeDir, err := GetWaveHomeDir()
	if err != nil {
		return "", err
	}
	logDir := filepath.Join(homeDir, "log")
	if err := ensureDir(logDir); err != nil {
		return "", fmt.Errorf("cannot create wave log directory %q: %w", logDir, err)
	}
	return logDir, nil
}

// GetOS returns the current operating system identifier.
func GetOS() string {
	return runtime.GOOS
}

// GetArch returns the current CPU architecture identifier.
func GetArch() string {
	return runtime.GOARCH
}

// ensureDir creates the directory at the given path if it does not already exist.
func ensureDir(path string) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(path, dirPermissions)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path %q exists but is not a directory", path)
	}
	return nil
}
