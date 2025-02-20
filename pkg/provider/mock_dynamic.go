package provider

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// MockDynamicClient implements dynamic.Interface for testing
type MockDynamicClient struct {
	dynamic.Interface
}

func (m *MockDynamicClient) Resource(resource schema.GroupVersionResource) dynamic.NamespaceableResourceInterface {
	return &MockNamespaceableResourceClient{}
}

type MockNamespaceableResourceClient struct {
	dynamic.NamespaceableResourceInterface
}

func (m *MockNamespaceableResourceClient) Create(ctx context.Context, obj *unstructured.Unstructured, opts metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (m *MockNamespaceableResourceClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions, subresources ...string) error {
	return nil
}

func (m *MockNamespaceableResourceClient) Get(ctx context.Context, name string, opts metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (m *MockNamespaceableResourceClient) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return &unstructured.UnstructuredList{}, nil
}
