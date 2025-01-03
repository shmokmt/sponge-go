package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestSponge_Run(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		append     bool
		existing   string
		wantOutput string
	}{
		{
			name:       "write to a new file",
			input:      "hello world",
			append:     false,
			existing:   "",
			wantOutput: "hello world",
		},
		{
			name:       "append to an existing file",
			input:      "hello world",
			append:     true,
			existing:   "existing content\n",
			wantOutput: "existing content\nhello world",
		},
		{
			name:       "overwrite an existing file",
			input:      "new content",
			append:     false,
			existing:   "old content",
			wantOutput: "new content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outFile, err := os.CreateTemp("", "sponge_test.*.txt")
			if err != nil {
				t.Fatalf("failed to create a temp file: %v", err)
			}
			defer os.Remove(outFile.Name())

			if tt.existing != "" {
				if _, err := outFile.WriteString(tt.existing); err != nil {
					t.Fatalf("failed to write existing content: %v", err)
				}
			}

			sponge := NewSponge(io.Discard, io.Discard, strings.NewReader(tt.input))

			err = sponge.Run(outFile.Name(), tt.append)
			if err != nil {
				t.Fatalf("Sponge.Run() error = %v, want nil", err)
			}

			output, err := os.ReadFile(outFile.Name())
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}

			if string(output) != tt.wantOutput {
				t.Errorf("Sponge.Run() output = %v, want %v", string(output), tt.wantOutput)
			}
		})
	}
}
