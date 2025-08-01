package main

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestParseFlagsDefaultValue(t *testing.T) {
	// Reset flag.CommandLine to avoid conflicts with other tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up args with no flags (should use default)
	os.Args = []string{"sift"}

	interval := parseFlags()
	expected := 3 * time.Second

	if interval != expected {
		t.Errorf("Expected default interval %v, got %v", expected, interval)
	}
}

func TestParseFlagsCustomValue(t *testing.T) {
	// Reset flag.CommandLine to avoid conflicts with other tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up args with custom refresh interval
	os.Args = []string{"sift", "--refresh-interval", "10"}

	interval := parseFlags()
	expected := 10 * time.Second

	if interval != expected {
		t.Errorf("Expected custom interval %v, got %v", expected, interval)
	}
}

func TestParseFlagsShortForm(t *testing.T) {
	// Reset flag.CommandLine to avoid conflicts with other tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up args with short form flag
	os.Args = []string{"sift", "--refresh-interval=5"}

	interval := parseFlags()
	expected := 5 * time.Second

	if interval != expected {
		t.Errorf("Expected interval %v, got %v", expected, interval)
	}
}

func TestParseFlagsZeroValue(t *testing.T) {
	// Reset flag.CommandLine to avoid conflicts with other tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up args with zero refresh interval
	os.Args = []string{"sift", "--refresh-interval", "0"}

	interval := parseFlags()
	expected := 0 * time.Second

	if interval != expected {
		t.Errorf("Expected zero interval %v, got %v", expected, interval)
	}
}

func TestParseFlagsLargeValue(t *testing.T) {
	// Reset flag.CommandLine to avoid conflicts with other tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up args with large refresh interval
	os.Args = []string{"sift", "--refresh-interval", "3600"}

	interval := parseFlags()
	expected := 3600 * time.Second

	if interval != expected {
		t.Errorf("Expected large interval %v, got %v", expected, interval)
	}
}
