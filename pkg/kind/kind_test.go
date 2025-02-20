package kind

// Variables for testing
var (
	ClusterExists = clusterExists
	CreateCluster = createCluster
	DeleteCluster = deleteCluster
)

func clusterExists(name string) (bool, error) {
	// Implementation
	return false, nil
}

func createCluster(configFile, name string) error {
	// Implementation
	return nil
}

func deleteCluster(name string) error {
	// Implementation
	return nil
}
