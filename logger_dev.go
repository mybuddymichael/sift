//go:build !prod

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getXDGCacheDir() (string, error) {
	if cacheDir := os.Getenv("XDG_CACHE_HOME"); cacheDir != "" {
		return cacheDir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache"), nil
}

var Logger *log.Logger

func init() {
	cacheDir, err := getXDGCacheDir()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(cacheDir, "sift")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	output := lumberjack.Logger{
		Filename:   filepath.Join(dir, "log"),
		MaxSize:    1,
		MaxBackups: 1,
	}
	Logger = log.NewWithOptions(&output, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           log.DebugLevel,
	})
}
