package crosslab

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/kanzifucius/crosslab/pkg/kind"
	"github.com/kanzifucius/crosslab/pkg/provider"
)

func TestInstallCrossplane(t *testing.T) {
	// Create a global test cluster name
	globalClusterName := fmt.Sprintf("test-cluster-%d", time.Now().Unix())

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful installation",
			wantErr: false,
		},
		{
			name:    "installation on existing crossplane",
			wantErr: false,
		},
	}

	// Create a real Kind client for testing
	kindClient := kind.NewClient()

	// Clean up any existing cluster with this name at the start
	_ = kindClient.DeleteCluster(globalClusterName)

	// Create initial cluster
	err := kindClient.CreateCluster(getKindConfig(), globalClusterName)
	if err != nil {
		t.Fatalf("failed to create test cluster: %v", err)
	}

	// Clean up after all tests
	defer func() {
		_ = kindClient.DeleteCluster(globalClusterName)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := provider.NewManager()
			if err != nil {
				t.Fatalf("failed to create manager: %v", err)
			}

			err = installCrossplane(context.Background(), manager)
			if (err != nil) != tt.wantErr {
				t.Errorf("installCrossplane() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallClusterProvider(t *testing.T) {
	// Create a global test cluster name
	globalClusterName := fmt.Sprintf("test-cluster-%d", time.Now().Unix())

	tests := []struct {
		name     string
		provider provider.Provider
		wantErr  bool
	}{
		{
			name: "successful provider installation",
			provider: provider.Provider{
				Name:    "provider-aws",
				Package: "crossplane/provider-aws",
				Version: "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "install same provider again",
			provider: provider.Provider{
				Name:    "provider-aws",
				Package: "crossplane/provider-aws",
				Version: "v1.0.0",
			},
			wantErr: false,
		},
	}

	// Create a real Kind client for testing
	kindClient := kind.NewClient()

	// Clean up any existing cluster with this name at the start
	_ = kindClient.DeleteCluster(globalClusterName)

	// Create initial cluster
	err := kindClient.CreateCluster("../../examples/kind-config.yaml", globalClusterName)
	if err != nil {
		t.Fatalf("failed to create test cluster: %v", err)
	}

	// Install Crossplane first
	manager, err := provider.NewManager()
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	err = installCrossplane(context.Background(), manager)
	if err != nil {
		t.Fatalf("failed to install Crossplane: %v", err)
	}

	// Clean up after all tests
	defer func() {
		_ = kindClient.DeleteCluster(globalClusterName)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := provider.NewManager()
			if err != nil {
				t.Fatalf("failed to create manager: %v", err)
			}

			err = installClusterProvider(context.Background(), manager, tt.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("installClusterProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getKindConfig() string {
	// Generate random ports in the range 30000-32767
	httpPort := 30000 + rand.Intn(2768)
	httpsPort := 30000 + rand.Intn(2768)

	// Ensure we don't get the same port for both
	for httpsPort == httpPort {
		httpsPort = 30000 + rand.Intn(2768)
	}

	kindConfig := fmt.Sprintf(`kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: %d
    protocol: TCP
  - containerPort: 443
    hostPort: %d
    protocol: TCP
- role: worker
- role: worker
networking:
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.96.0.0/16"`, httpPort, httpsPort)

	return kindConfig
}

func TestCreateCmd(t *testing.T) {
	// Create a global test cluster name
	globalClusterName := fmt.Sprintf("test-cluster-%d", time.Now().Unix())

	// Create temporary KIND config file
	kindConfig := getKindConfig()

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
		name        string
		configFile  string
		forceCreate bool
		wantErr     bool
	}{
		{
			name:       "successful cluster creation",
			configFile: tmpfile.Name(),
			wantErr:    false,
		},
		{
			name:       "cluster already exists",
			configFile: tmpfile.Name(),
			wantErr:    true,
		},
		{
			name:        "force recreation of existing cluster",
			configFile:  tmpfile.Name(),
			forceCreate: true,
			wantErr:     false,
		},
	}

	// Create a real Kind client for testing
	kindClient := kind.NewClient()

	// Clean up any existing cluster with this name at the start
	_ = kindClient.DeleteCluster(globalClusterName)

	// Clean up after all tests
	defer func() {
		_ = kindClient.DeleteCluster(globalClusterName)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the "cluster already exists" test, create a cluster first
			if tt.name == "cluster already exists" {
				err := kindClient.CreateCluster(tt.configFile, globalClusterName)
				if err != nil {
					t.Fatalf("failed to create test cluster: %v", err)
				}
			}

			// Set up command flags
			kindConfigFile = tt.configFile
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
