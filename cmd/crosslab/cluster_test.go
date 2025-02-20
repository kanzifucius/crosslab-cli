package crosslab

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kanzifucius/crosslab/pkg/config"
	"github.com/kanzifucius/crosslab/pkg/kind"
	"github.com/kanzifucius/crosslab/pkg/provider"
	"gopkg.in/yaml.v3"
)

// cleanupTestClusters deletes all KIND clusters that match the test prefix pattern
func cleanupTestClusters(t *testing.T) {
	t.Helper()
	kindManager := kind.NewManager()

	// Get list of all clusters
	clusters, err := kindManager.ListClusters()
	if err != nil {
		t.Logf("Warning: failed to list clusters for cleanup: %v", err)
		return
	}

	// Delete any test clusters (those with "test-cluster-" prefix)
	for _, cluster := range clusters {
		if strings.HasPrefix(cluster, "test-cluster-") {
			t.Logf("Cleaning up existing test cluster: %s", cluster)
			if err := kindManager.DeleteCluster(cluster); err != nil {
				t.Logf("Warning: failed to delete test cluster %s: %v", cluster, err)
			}
		}
	}
}

func TestInstallCrossplane(t *testing.T) {
	// Clean up any existing test clusters first
	cleanupTestClusters(t)

	// Create a global test cluster name
	globalClusterName := fmt.Sprintf("test-cluster-%d", time.Now().Unix())

	// Create temporary KIND config file
	kindConfigObj := config.DefaultKindConfig()
	kindConfigBytes, err := yaml.Marshal(kindConfigObj)
	if err != nil {
		t.Fatalf("failed to marshal KIND config: %v", err)
	}
	kindConfig := string(kindConfigBytes)

	tmpfile, err := os.CreateTemp("", "kind-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(kindConfig)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful installation",
			wantErr: false,
		},
	}

	// Create a real Kind manager for testing
	kindManager := kind.NewManager()

	// Clean up any existing cluster with this name at the start
	_ = kindManager.DeleteCluster(globalClusterName)

	// Create initial cluster
	err = kindManager.CreateCluster(tmpfile.Name(), globalClusterName)
	if err != nil {
		t.Fatalf("failed to create test cluster: %v", err)
	}

	// Clean up after all tests
	defer func() {
		_ = kindManager.DeleteCluster(globalClusterName)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//delete all test clusters
			cleanupTestClusters(t)
			manager, err := provider.NewManager()
			if err != nil {
				t.Fatalf("failed to create manager: %v", err)
			}

			err = InstallCrossplane(context.Background(), manager)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallCrossplane() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallClusterProvider(t *testing.T) {
	// Clean up any existing test clusters first
	cleanupTestClusters(t)

	// Create a global test cluster name
	globalClusterName := fmt.Sprintf("test-cluster-%d", time.Now().Unix())

	// Create temporary KIND config file
	kindConfigObj := config.DefaultKindConfig()
	kindConfigBytes, err := yaml.Marshal(kindConfigObj)
	if err != nil {
		t.Fatalf("failed to marshal KIND config: %v", err)
	}
	kindConfig := string(kindConfigBytes)

	tmpfile, err := os.CreateTemp("", "kind-config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(kindConfig)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	tests := []struct {
		name     string
		provider config.Provider
		wantErr  bool
	}{
		{
			name: "successful provider installation",
			provider: config.Provider{
				Name:    "provider-aws",
				Package: "crossplane/provider-aws",
				Version: "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "install same provider again",
			provider: config.Provider{
				Name:    "provider-aws",
				Package: "crossplane/provider-aws",
				Version: "v1.0.0",
			},
			wantErr: false,
		},
	}

	// Create a real Kind manager for testing
	kindManager := kind.NewManager()

	// Clean up any existing cluster with this name at the start
	_ = kindManager.DeleteCluster(globalClusterName)

	// Create initial cluster
	err = kindManager.CreateCluster(tmpfile.Name(), globalClusterName)
	if err != nil {
		t.Fatalf("failed to create test cluster: %v", err)
	}

	// Install Crossplane first
	manager, err := provider.NewManager()
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	err = InstallCrossplane(context.Background(), manager)
	if err != nil {
		t.Fatalf("failed to install Crossplane: %v", err)
	}

	// Clean up after all tests
	defer func() {
		_ = kindManager.DeleteCluster(globalClusterName)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := provider.NewManager()
			if err != nil {
				t.Fatalf("failed to create manager: %v", err)
			}

			err = InstallClusterProvider(context.Background(), manager, tt.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallClusterProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCmd(t *testing.T) {
	// Clean up any existing test clusters first
	cleanupTestClusters(t)

	//create temp dir
	tmpDir, err := os.MkdirTemp("", "crosslab-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	// generate random cluster name
	globalClusterName := fmt.Sprintf("test-cluster-%d", rand.Intn(1000000))

	// Create temporary KIND config file
	initializer := config.NewInitializer(tmpDir + "/.crosslab")
	err = initializer.Initialize()
	if err != nil {
		t.Fatalf("failed to initialize config: %v", err)
	}
	kindConfigFile = initializer.GetKindConfig()
	clusterConfig = initializer.GetProvidersConfig()

	tests := []struct {
		name        string
		forceCreate bool
		wantErr     bool
	}{
		{
			name:    "successful cluster creation",
			wantErr: false,
		},
	}

	// Create a real Kind manager for testing
	kindManager := kind.NewManager()

	// Clean up after all tests
	defer func() {
		_ = kindManager.DeleteCluster(globalClusterName)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterName = globalClusterName
			forceCreate = tt.forceCreate

			// Execute command
			err := createCmd.RunE(createCmd, []string{})
			if (err != nil) != tt.wantErr {
				t.Errorf("createCmd.RunE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
