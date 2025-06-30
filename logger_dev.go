//go:build log

package main

import (
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.DebugLevel.String())).
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("4"))
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.InfoLevel.String())).
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("7"))
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Level:           log.DebugLevel,
	})
	Logger.SetStyles(styles)
}
