//go:build !prod

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *log.Logger

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(home, ".prioritizer-terminal")
	if err := os.MkdirAll(dir, 0755); err != nil {
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
		TimeFormat:      time.Kitchen,
		Level:           log.DebugLevel,
	})
}
