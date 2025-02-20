package provider

import (
	"context"

	"github.com/kanzifucius/crosslab/pkg/config"
)

// mockManager implements provider operations for testing
type mockManager struct {
	InstallCrossplaneFunc func(ctx context.Context) error
	WaitForCrossplaneFunc func(ctx context.Context) error
	InstallFunc           func(ctx context.Context, p config.Provider, force bool) error
	WaitForHealthFunc     func(ctx context.Context, name string) error
	ListFunc              func(ctx context.Context) ([]config.Provider, error)
	DeleteFunc            func(ctx context.Context, name string) error
	ExistsFunc            func(ctx context.Context, name string) (bool, error)
}

// NewMockManager creates a new mock provider manager
func NewMockManager() Manager {
	return &mockManager{}
}

func (m *mockManager) InstallCrossplane(ctx context.Context) error {
	if m.InstallCrossplaneFunc != nil {
		return m.InstallCrossplaneFunc(ctx)
	}
	return nil
}

func (m *mockManager) WaitForCrossplaneHealth(ctx context.Context) error {
	if m.WaitForCrossplaneFunc != nil {
		return m.WaitForCrossplaneFunc(ctx)
	}
	return nil
}

func (m *mockManager) Install(ctx context.Context, p config.Provider, force bool) error {
	if m.InstallFunc != nil {
		return m.InstallFunc(ctx, p, force)
	}
	return nil
}

func (m *mockManager) WaitForHealth(ctx context.Context, name string) error {
	if m.WaitForHealthFunc != nil {
		return m.WaitForHealthFunc(ctx, name)
	}
	return nil
}

func (m *mockManager) List(ctx context.Context) ([]config.Provider, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *mockManager) Delete(ctx context.Context, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, name)
	}
	return nil
}

func (m *mockManager) Exists(ctx context.Context, name string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, name)
	}
	return false, nil
}
