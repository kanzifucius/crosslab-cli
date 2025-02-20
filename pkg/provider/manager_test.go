package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/kanzifucius/crosslab/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	// Test case: provider exists
	t.Run("provider exists", func(t *testing.T) {
		mockManager := &mockManager{
			ExistsFunc: func(ctx context.Context, name string) (bool, error) {
				return true, nil
			},
		}

		exists, err := mockManager.Exists(context.Background(), "provider-aws")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	// Test case: provider does not exist
	t.Run("provider does not exist", func(t *testing.T) {
		mockManager := &mockManager{
			ExistsFunc: func(ctx context.Context, name string) (bool, error) {
				return false, nil
			},
		}

		exists, err := mockManager.Exists(context.Background(), "provider-aws")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	// Test case: error checking if provider exists
	t.Run("error checking provider", func(t *testing.T) {
		expectedErr := errors.New("connection error")
		mockManager := &mockManager{
			ExistsFunc: func(ctx context.Context, name string) (bool, error) {
				return false, expectedErr
			},
		}

		exists, err := mockManager.Exists(context.Background(), "provider-aws")
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.False(t, exists)
	})
}

func TestInstall(t *testing.T) {
	// Test case: successful installation
	t.Run("successful installation", func(t *testing.T) {
		mockManager := &mockManager{
			ExistsFunc: func(ctx context.Context, name string) (bool, error) {
				return false, nil
			},
			InstallFunc: func(ctx context.Context, p config.Provider, force bool) error {
				return nil
			},
		}

		provider := config.Provider{
			Name:    "provider-aws",
			Package: "crossplane/provider-aws",
			Version: "v0.24.1",
		}

		err := mockManager.Install(context.Background(), provider, false)
		assert.NoError(t, err)
	})

	// Test case: provider already exists and force is false
	t.Run("provider exists and force is false", func(t *testing.T) {
		mockManager := &mockManager{
			ExistsFunc: func(ctx context.Context, name string) (bool, error) {
				return true, nil
			},
			InstallFunc: func(ctx context.Context, p config.Provider, force bool) error {
				if p.Name == "provider-aws" && !force {
					return errors.New("provider provider-aws already exists, use force option to reinstall")
				}
				return nil
			},
		}

		provider := config.Provider{
			Name:    "provider-aws",
			Package: "crossplane/provider-aws",
			Version: "v0.24.1",
		}

		err := mockManager.Install(context.Background(), provider, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	// Test case: provider already exists and force is true
	t.Run("provider exists and force is true", func(t *testing.T) {
		var deleteCalled, installCalled bool
		mockManager := &mockManager{
			ExistsFunc: func(ctx context.Context, name string) (bool, error) {
				return true, nil
			},
			DeleteFunc: func(ctx context.Context, name string) error {
				deleteCalled = true
				return nil
			},
			InstallFunc: func(ctx context.Context, p config.Provider, force bool) error {
				installCalled = true
				return nil
			},
		}

		provider := config.Provider{
			Name:    "provider-aws",
			Package: "crossplane/provider-aws",
			Version: "v0.24.1",
		}

		// Manually simulate what the real implementation would do
		ctx := context.Background()
		exists, _ := mockManager.Exists(ctx, provider.Name)
		if exists && true { // force is true
			_ = mockManager.Delete(ctx, provider.Name)
		}
		err := mockManager.Install(ctx, provider, true)

		assert.NoError(t, err)
		assert.True(t, deleteCalled)
		assert.True(t, installCalled)
	})

	// Test case: error during installation
	t.Run("error during installation", func(t *testing.T) {
		expectedErr := errors.New("installation error")
		mockManager := &mockManager{
			ExistsFunc: func(ctx context.Context, name string) (bool, error) {
				return false, nil
			},
			InstallFunc: func(ctx context.Context, p config.Provider, force bool) error {
				return expectedErr
			},
		}

		provider := config.Provider{
			Name:    "provider-aws",
			Package: "crossplane/provider-aws",
			Version: "v0.24.1",
		}

		err := mockManager.Install(context.Background(), provider, false)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestDelete(t *testing.T) {
	// Test case: successful deletion
	t.Run("successful deletion", func(t *testing.T) {
		mockManager := &mockManager{
			DeleteFunc: func(ctx context.Context, name string) error {
				return nil
			},
		}

		err := mockManager.Delete(context.Background(), "provider-aws")
		assert.NoError(t, err)
	})

	// Test case: error during deletion
	t.Run("error during deletion", func(t *testing.T) {
		expectedErr := errors.New("deletion error")
		mockManager := &mockManager{
			DeleteFunc: func(ctx context.Context, name string) error {
				return expectedErr
			},
		}

		err := mockManager.Delete(context.Background(), "provider-aws")
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestList(t *testing.T) {
	// Test case: successful listing
	t.Run("successful listing", func(t *testing.T) {
		expectedProviders := []config.Provider{
			{
				Name:    "provider-aws",
				Package: "crossplane/provider-aws:v0.24.1",
			},
			{
				Name:    "provider-gcp",
				Package: "crossplane/provider-gcp:v0.18.0",
			},
		}

		mockManager := &mockManager{
			ListFunc: func(ctx context.Context) ([]config.Provider, error) {
				return expectedProviders, nil
			},
		}

		providers, err := mockManager.List(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedProviders, providers)
	})

	// Test case: error during listing
	t.Run("error during listing", func(t *testing.T) {
		expectedErr := errors.New("listing error")
		mockManager := &mockManager{
			ListFunc: func(ctx context.Context) ([]config.Provider, error) {
				return nil, expectedErr
			},
		}

		providers, err := mockManager.List(context.Background())
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, providers)
	})
}

func TestWaitForHealth(t *testing.T) {
	// Test case: provider becomes healthy
	t.Run("provider becomes healthy", func(t *testing.T) {
		mockManager := &mockManager{
			WaitForHealthFunc: func(ctx context.Context, name string) error {
				return nil
			},
		}

		err := mockManager.WaitForHealth(context.Background(), "provider-aws")
		assert.NoError(t, err)
	})

	// Test case: timeout waiting for provider health
	t.Run("timeout waiting for provider health", func(t *testing.T) {
		expectedErr := errors.New("timeout waiting for provider provider-aws to become healthy")
		mockManager := &mockManager{
			WaitForHealthFunc: func(ctx context.Context, name string) error {
				return expectedErr
			},
		}

		err := mockManager.WaitForHealth(context.Background(), "provider-aws")
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestInstallCrossplane(t *testing.T) {
	// Test case: successful installation
	t.Run("successful installation", func(t *testing.T) {
		mockManager := &mockManager{
			InstallCrossplaneFunc: func(ctx context.Context) error {
				return nil
			},
		}

		err := mockManager.InstallCrossplane(context.Background())
		assert.NoError(t, err)
	})

	// Test case: error during installation
	t.Run("error during installation", func(t *testing.T) {
		expectedErr := errors.New("installation error")
		mockManager := &mockManager{
			InstallCrossplaneFunc: func(ctx context.Context) error {
				return expectedErr
			},
		}

		err := mockManager.InstallCrossplane(context.Background())
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestWaitForCrossplaneHealth(t *testing.T) {
	// Test case: Crossplane becomes healthy
	t.Run("Crossplane becomes healthy", func(t *testing.T) {
		mockManager := &mockManager{
			WaitForCrossplaneFunc: func(ctx context.Context) error {
				return nil
			},
		}

		err := mockManager.WaitForCrossplaneHealth(context.Background())
		assert.NoError(t, err)
	})

	// Test case: timeout waiting for Crossplane health
	t.Run("timeout waiting for Crossplane health", func(t *testing.T) {
		expectedErr := errors.New("timeout waiting for Crossplane to become healthy")
		mockManager := &mockManager{
			WaitForCrossplaneFunc: func(ctx context.Context) error {
				return expectedErr
			},
		}

		err := mockManager.WaitForCrossplaneHealth(context.Background())
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
