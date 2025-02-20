package provider

import "errors"

var (
	// ErrInstallFailed indicates that the installation of a provider or Crossplane failed
	ErrInstallFailed = errors.New("installation failed")

	// ErrHealthCheckFailed indicates that the health check for a provider or Crossplane failed
	ErrHealthCheckFailed = errors.New("health check failed")
)
