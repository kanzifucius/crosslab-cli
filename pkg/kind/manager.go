package kind

import (
	"fmt"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

// Manager defines the operations that can be performed on a Kind cluster
type Manager interface {
	// ClusterExists checks if a cluster with the given name exists
	ClusterExists(name string) (bool, error)
	// CreateCluster creates a new Kind cluster using the provided configuration file
	CreateCluster(configFilePath string, name string) error
	// DeleteCluster deletes a Kind cluster by name
	DeleteCluster(name string) error
	// ListClusters returns a list of existing Kind clusters
	ListClusters() ([]string, error)
}

// manager handles Kind cluster operations
type manager struct {
	provider cluster.Provider
}

// NewManager creates a new Kind cluster manager
func NewManager() Manager {
	return &manager{
		provider: *cluster.NewProvider(
			cluster.ProviderWithLogger(cmd.NewLogger()),
		),
	}
}

// ClusterExists checks if a cluster with the given name exists
func (m *manager) ClusterExists(name string) (bool, error) {
	clusters, err := m.provider.List()
	if err != nil {
		return false, fmt.Errorf("error checking cluster existence: %v", err)
	}

	for _, cluster := range clusters {
		if cluster == name {
			return true, nil
		}
	}

	return false, nil
}

// CreateCluster creates a new Kind cluster using the provided configuration file
func (m *manager) CreateCluster(configFilePath string, name string) error {
	exists, err := m.ClusterExists(name)
	if err != nil {
		return fmt.Errorf("error checking cluster existence: %v", err)
	}

	if exists {
		return fmt.Errorf("cluster %s already exists", name)
	}

	// Create the cluster
	if err := m.provider.Create(
		name,
		cluster.CreateWithConfigFile(configFilePath),
	); err != nil {
		return fmt.Errorf("error creating cluster: %v", err)
	}

	return nil
}

// DeleteCluster deletes a Kind cluster by name
func (m *manager) DeleteCluster(name string) error {
	if err := m.provider.Delete(name, ""); err != nil {
		return fmt.Errorf("error deleting cluster: %v", err)
	}

	return nil
}

// ListClusters returns a list of existing Kind clusters
func (m *manager) ListClusters() ([]string, error) {
	clusters, err := m.provider.List()
	if err != nil {
		return nil, fmt.Errorf("error listing clusters: %v", err)
	}

	return clusters, nil
}
