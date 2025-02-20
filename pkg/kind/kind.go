package kind

import (
	"fmt"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

// Interface defines the operations that can be performed on a Kind cluster
type Interface interface {
	ClusterExists(name string) (bool, error)
	CreateCluster(configFilePath string, name string) error
	DeleteCluster(name string) error
	ListClusters() ([]string, error)
}

// Client handles Kind cluster operations
type Client struct {
	provider cluster.Provider
}

// For testing
var NewClient = newClient

// newClient creates a new Kind client
func newClient() Interface {
	return &Client{
		provider: *cluster.NewProvider(
			cluster.ProviderWithLogger(cmd.NewLogger()),
		),
	}
}

// ClusterExists checks if a cluster with the given name exists
func (c *Client) ClusterExists(name string) (bool, error) {
	clusters, err := c.provider.List()
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
func (c *Client) CreateCluster(configFilePath string, name string) error {
	exists, err := c.ClusterExists(name)
	if err != nil {
		return fmt.Errorf("error checking cluster existence: %v", err)
	}

	if exists {
		return fmt.Errorf("cluster %s already exists", name)
	}

	// Create the cluster
	if err := c.provider.Create(
		name,
		cluster.CreateWithConfigFile(configFilePath),
	); err != nil {
		return fmt.Errorf("error creating cluster: %v", err)
	}

	return nil
}

// DeleteCluster deletes a Kind cluster by name
func (c *Client) DeleteCluster(name string) error {
	if err := c.provider.Delete(name, ""); err != nil {
		return fmt.Errorf("error deleting cluster: %v", err)
	}

	return nil
}

// ListClusters returns a list of existing Kind clusters
func (c *Client) ListClusters() ([]string, error) {
	clusters, err := c.provider.List()
	if err != nil {
		return nil, fmt.Errorf("error listing clusters: %v", err)
	}

	return clusters, nil
}
