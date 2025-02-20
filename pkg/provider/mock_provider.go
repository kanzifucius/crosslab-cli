package provider

import (
	"context"
)

// MockManager implements provider operations for testing
type MockManager struct {
	InstallCrossplaneFunc func(ctx context.Context) error
	WaitForCrossplaneFunc func(ctx context.Context) error
	InstallFunc           func(ctx context.Context, p Provider, force bool) error
	WaitForHealthFunc     func(ctx context.Context, name string) error
	ListFunc              func(ctx context.Context) ([]Provider, error)
	DeleteFunc            func(ctx context.Context, name string) error
	ExistsFunc            func(ctx context.Context, name string) (bool, error)
}

func (m *MockManager) InstallCrossplane(ctx context.Context) error {
	if m.InstallCrossplaneFunc != nil {
		return m.InstallCrossplaneFunc(ctx)
	}
	return nil
}

func (m *MockManager) WaitForCrossplaneHealth(ctx context.Context) error {
	if m.WaitForCrossplaneFunc != nil {
		return m.WaitForCrossplaneFunc(ctx)
	}
	return nil
}

func (m *MockManager) Install(ctx context.Context, p Provider, force bool) error {
	if m.InstallFunc != nil {
		return m.InstallFunc(ctx, p, force)
	}
	return nil
}

func (m *MockManager) WaitForHealth(ctx context.Context, name string) error {
	if m.WaitForHealthFunc != nil {
		return m.WaitForHealthFunc(ctx, name)
	}
	return nil
}

func (m *MockManager) List(ctx context.Context) ([]Provider, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockManager) Delete(ctx context.Context, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, name)
	}
	return nil
}

func (m *MockManager) Exists(ctx context.Context, name string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, name)
	}
	return false, nil
}
