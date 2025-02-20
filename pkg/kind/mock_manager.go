package kind

// mockManager implements Manager interface for testing
type mockManager struct {
	ClusterExistsFunc func(name string) (bool, error)
	CreateClusterFunc func(configFilePath string, name string) error
	DeleteClusterFunc func(name string) error
	ListClustersFunc  func() ([]string, error)
}

// NewMockManager creates a new mock kind cluster manager
func NewMockManager() Manager {
	return &mockManager{}
}

func (m *mockManager) ClusterExists(name string) (bool, error) {
	if m.ClusterExistsFunc != nil {
		return m.ClusterExistsFunc(name)
	}
	return false, nil
}

func (m *mockManager) CreateCluster(configFilePath string, name string) error {
	if m.CreateClusterFunc != nil {
		return m.CreateClusterFunc(configFilePath, name)
	}
	return nil
}

func (m *mockManager) DeleteCluster(name string) error {
	if m.DeleteClusterFunc != nil {
		return m.DeleteClusterFunc(name)
	}
	return nil
}

func (m *mockManager) ListClusters() ([]string, error) {
	if m.ListClustersFunc != nil {
		return m.ListClustersFunc()
	}
	return nil, nil
}
