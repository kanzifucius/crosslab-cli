package config

import (
	kindconfigv1alpha4 "sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
)

// DefaultKindConfig returns a default Kind cluster configuration
func DefaultKindConfig() *kindconfigv1alpha4.Cluster {
	return &kindconfigv1alpha4.Cluster{
		TypeMeta: kindconfigv1alpha4.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "kind.x-k8s.io/v1alpha4",
		},
		Nodes: []kindconfigv1alpha4.Node{
			{
				Role: kindconfigv1alpha4.ControlPlaneRole,
				ExtraPortMappings: []kindconfigv1alpha4.PortMapping{
					{
						ContainerPort: 80,
						HostPort:      8080,
						Protocol:      "TCP",
					},
					{
						ContainerPort: 443,
						HostPort:      8443,
						Protocol:      "TCP",
					},
				},
			},
			{
				Role: kindconfigv1alpha4.WorkerRole,
			},
			{
				Role: kindconfigv1alpha4.WorkerRole,
			},
		},
		Networking: kindconfigv1alpha4.Networking{
			PodSubnet:     "10.244.0.0/16",
			ServiceSubnet: "10.96.0.0/16",
		},
	}
}
