package crosslab

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kanzifucius/crosslab/pkg/config"
	"github.com/spf13/cobra"
)

func TestInitCmd(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "crosslab-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	tests := []struct {
		name       string
		outputDir  string
		wantErr    bool
		checkFiles []string
	}{
		{
			name:      "default initialization",
			outputDir: filepath.Join(tempDir, ".crosslab"),
			wantErr:   false,
			checkFiles: []string{
				"kind-config.yaml",
				"config/crosslab-config.yaml",
			},
		},
		{
			name:      "custom directory",
			outputDir: filepath.Join(tempDir, "custom-dir"),
			wantErr:   false,
			checkFiles: []string{
				"kind-config.yaml",
				"config/crosslab-config.yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command for each test to avoid flag conflicts
			cmd := &cobra.Command{
				Use: "init",
			}

			// Set the output directory flag
			cmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".crosslab", "Output directory")
			outputDir = tt.outputDir

			// Execute the command
			err := initCmd.RunE(cmd, []string{})

			// Check error status
			if (err != nil) != tt.wantErr {
				t.Errorf("initCmd.RunE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error expected, check that files were created
			if !tt.wantErr {
				for _, file := range tt.checkFiles {
					fullPath := filepath.Join(tt.outputDir, file)
					if _, err := os.Stat(fullPath); os.IsNotExist(err) {
						t.Errorf("expected file %s does not exist", fullPath)
					}
				}

				// Verify the config file structure
				configPath := filepath.Join(tt.outputDir, "config/crosslab-config.yaml")
				cfg, err := config.LoadConfig(configPath)
				if err != nil {
					t.Errorf("failed to load config file: %v", err)
					return
				}

				// Validate the config structure
				if err := cfg.Validate(); err != nil {
					t.Errorf("config validation failed: %v", err)
				}
			}
		})
	}
}
