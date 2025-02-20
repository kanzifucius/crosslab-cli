package provider

import (
	"context"
	"fmt"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	CrossplaneNamespace = "crossplane-system"
	ProviderTimeout     = 300 * time.Second
	CrossplaneHelmRepo  = "https://charts.crossplane.io/stable"
	CrossplaneChartName = "crossplane"
)

// Manager handles Crossplane provider operations
type Manager struct {
	Client dynamic.Interface
}

// For testing
var NewManager = newManager

// newManager creates a new provider manager
func newManager() (*Manager, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	return &Manager{Client: client}, nil
}

// getKubeConfig returns a Kubernetes REST config
func getKubeConfig() (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes config: %v", err)
	}

	return config, nil
}

// Delete deletes a Crossplane provider
func (m *Manager) Delete(ctx context.Context, providerName string) error {
	providerGVR := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	if err := m.Client.Resource(providerGVR).Delete(ctx, providerName, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete provider %s: %v", providerName, err)
	}

	return nil
}

// Exists checks if a provider already exists
func (m *Manager) Exists(ctx context.Context, providerName string) (bool, error) {
	providerGVR := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	_, err := m.Client.Resource(providerGVR).Get(ctx, providerName, metav1.GetOptions{})
	if err != nil {
		return false, nil
	}

	return true, nil
}

// Install installs or updates a Crossplane provider
func (m *Manager) Install(ctx context.Context, provider Provider, force bool) error {
	exists, err := m.Exists(ctx, provider.Name)
	if err != nil {
		return fmt.Errorf("failed to check if provider exists: %v", err)
	}

	if exists {
		if !force {
			return fmt.Errorf("provider %s already exists, use force option to reinstall", provider.Name)
		}

		// Delete existing provider
		if err := m.Delete(ctx, provider.Name); err != nil {
			return fmt.Errorf("failed to delete existing provider: %v", err)
		}

		// Wait for provider to be deleted
		for {
			exists, err := m.Exists(ctx, provider.Name)
			if err != nil {
				return fmt.Errorf("failed to check if provider is deleted: %v", err)
			}
			if !exists {
				break
			}
			time.Sleep(2 * time.Second)
		}
	}

	providerGVR := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	providerObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "pkg.crossplane.io/v1",
			"kind":       "Provider",
			"metadata": map[string]interface{}{
				"name": provider.Name,
			},
			"spec": map[string]interface{}{
				"package": fmt.Sprintf("%s:%s", provider.Package, provider.Version),
			},
		},
	}

	_, err = m.Client.Resource(providerGVR).Create(ctx, providerObj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create provider %s: %v", provider.Name, err)
	}

	return nil
}

// WaitForHealth waits for a provider to become healthy
func (m *Manager) WaitForHealth(ctx context.Context, providerName string) error {
	providerGVR := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, ProviderTimeout)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for provider %s to become healthy", providerName)
		default:
			provider, err := m.Client.Resource(providerGVR).Get(ctx, providerName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get provider %s: %v", providerName, err)
			}

			conditions, found, err := unstructured.NestedSlice(provider.Object, "status", "conditions")
			if err != nil || !found {
				continue
			}

			for _, c := range conditions {
				condition, ok := c.(map[string]interface{})
				if !ok {
					continue
				}

				if condition["type"] == "Healthy" && condition["status"] == "True" {
					return nil
				}
			}

			time.Sleep(5 * time.Second)
		}
	}
}

// List returns a list of installed Crossplane providers
func (m *Manager) List(ctx context.Context) ([]Provider, error) {
	providerGVR := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	list, err := m.Client.Resource(providerGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %v", err)
	}

	var providers []Provider
	for _, item := range list.Items {
		pkg, found, err := unstructured.NestedString(item.Object, "spec", "package")
		if err != nil || !found {
			continue
		}

		providers = append(providers, Provider{
			Name:    item.GetName(),
			Package: pkg,
		})
	}

	return providers, nil
}

// InstallCrossplane installs the Crossplane Helm chart
func (m *Manager) InstallCrossplane(ctx context.Context) error {
	// Create namespace if it doesn't exist
	ns := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				"name": CrossplaneNamespace,
			},
		},
	}

	nsGVR := schema.GroupVersionResource{Version: "v1", Resource: "namespaces"}
	_, err := m.Client.Resource(nsGVR).Create(ctx, ns, metav1.CreateOptions{})
	if err != nil && !isAlreadyExists(err) {
		return fmt.Errorf("failed to create namespace: %v", err)
	}

	// Initialize Helm settings
	settings := cli.New()
	settings.SetNamespace(CrossplaneNamespace)

	// Add Crossplane Helm repository
	repoEntry := repo.Entry{
		Name: "crossplane-stable",
		URL:  CrossplaneHelmRepo,
	}

	// Create repository with HTTP client
	r, err := repo.NewChartRepository(&repoEntry, getter.All(settings))
	if err != nil {
		return fmt.Errorf("failed to create chart repository: %v", err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return fmt.Errorf("failed to download repository index: %v", err)
	}

	// Initialize Helm action configuration
	actionConfig := new(action.Configuration)
	err = actionConfig.Init(settings.RESTClientGetter(), CrossplaneNamespace, "secret", func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	})
	if err != nil {
		return fmt.Errorf("failed to initialize helm configuration: %v", err)
	}

	// Create Helm install client
	client := action.NewInstall(actionConfig)
	client.Namespace = CrossplaneNamespace
	client.CreateNamespace = true
	client.Wait = true
	client.Timeout = 5 * time.Minute
	client.ReleaseName = CrossplaneChartName
	client.ChartPathOptions.RepoURL = CrossplaneHelmRepo

	// Load Crossplane chart
	chartPath, err := client.ChartPathOptions.LocateChart("crossplane", settings)
	if err != nil {
		return fmt.Errorf("failed to locate Crossplane chart: %v", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load Crossplane chart: %v", err)
	}

	// Install Crossplane
	_, err = client.Run(chart, nil)
	if err != nil {
		return fmt.Errorf("failed to install Crossplane: %v", err)
	}

	return nil
}

// WaitForCrossplaneHealth waits for Crossplane to become healthy
func (m *Manager) WaitForCrossplaneHealth(ctx context.Context) error {
	deployGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, ProviderTimeout)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for Crossplane to become healthy")
		default:
			deploy, err := m.Client.Resource(deployGVR).Namespace(CrossplaneNamespace).
				Get(ctx, "crossplane", metav1.GetOptions{})
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			conditions, found, err := unstructured.NestedSlice(deploy.Object, "status", "conditions")
			if err != nil || !found {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, c := range conditions {
				condition, ok := c.(map[string]interface{})
				if !ok {
					continue
				}

				if condition["type"] == "Available" && condition["status"] == "True" {
					return nil
				}
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func isAlreadyExists(err error) bool {
	return err != nil && err.Error() == "already exists"
}
