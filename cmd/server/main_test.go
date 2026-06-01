package main

import (
	"io"
	"strings"
	"testing"
)

func TestParsePortDefaultsTo8080(t *testing.T) {
	got, err := parsePort(nil, io.Discard)
	if err != nil {
		t.Fatalf("parsePort() error = %v", err)
	}

	if got != defaultPort {
		t.Fatalf("port = %d, want %d", got, defaultPort)
	}
}

func TestParsePortReadsCommandLineFlag(t *testing.T) {
	got, err := parsePort([]string{"-port", "9090"}, io.Discard)
	if err != nil {
		t.Fatalf("parsePort() error = %v", err)
	}

	if got != 9090 {
		t.Fatalf("port = %d, want 9090", got)
	}
}

func TestParsePortReadsPositionalPortArgument(t *testing.T) {
	got, err := parsePort([]string{"9091"}, io.Discard)
	if err != nil {
		t.Fatalf("parsePort() error = %v", err)
	}

	if got != 9091 {
		t.Fatalf("port = %d, want 9091", got)
	}
}

func TestParsePortRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "zero", args: []string{"-port", "0"}},
		{name: "too high", args: []string{"-port", "65536"}},
		{name: "not a number", args: []string{"-port", "abc"}},
		{name: "not a positional number", args: []string{"abc"}},
		{name: "flag and positional", args: []string{"-port", "9090", "9091"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parsePort(tt.args, io.Discard)
			if err == nil {
				t.Fatal("parsePort() error = nil, want error")
			}
			if tt.name != "not a number" && !strings.Contains(err.Error(), "port") {
				t.Fatalf("parsePort() error = %q, want port validation error", err.Error())
			}
		})
	}
}
