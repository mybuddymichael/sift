//go:build prod

package main

import (
	"io"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	Logger = log.New(io.Discard)
}
