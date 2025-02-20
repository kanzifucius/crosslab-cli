# Crosslab

A Cobra-based CLI application for managing Kind clusters and Crossplane providers. This project enables the creation of a local development environment for Crossplane-based work such as composition creation, testing, and development.

This should not require the installation of Crossplane or Crossplane's dependencies nor kind.
However, it is typical to have these installed on your local machine.

## Table of Contents
- [Installation](#installation)
- [Getting Started](#getting-started)
- [CLI Commands](#cli-commands)
- [Kind Cluster Management](#kind-cluster-management)
- [Crossplane Provider Management](#crossplane-provider-management)
- [Development](#development)
- [Project Structure](#project-structure)

## Installation

### Prerequisites

- Go 1.23.2 or later
- Docker (for Kind clusters)
- Make
- GoReleaser (optional, for releases)

### Pre-built binaries

Pre-built binaries are available for download from the [releases page](https://github.com/kanzifucius/crosslab-cli/releases).

### Install script

The easiest way to install Crosslab is using the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/kanzifucius/crosslab-cli/main/install.sh | sh
```

### Building from source

To build from source:

```bash
git clone https://github.com/kanzifucius/crosslab-cli.git
cd crosslab-cli
make build
```

The binary will be available in the `dist` directory. MacOS and Linux users can use the `crosslab` binary directly, while Windows users will need to use the `crosslab.exe` binary.

## Getting Started

After installation, you can initialize Crosslab with:

```bash
crosslab init
```

This will create the necessary configuration files in the current directory (or specify a different directory with `--output-dir`).

## CLI Commands

### Core Commands

- `crosslab` - Root command
- `crosslab version` - Show the CLI version
- `crosslab init` - Initialize configuration files
  - `--output-dir, -o` - Output directory for configuration files (default: current directory)

### Cluster Management

- `crosslab cluster` - Manage Kind clusters
  - `create` - Create a new Kind cluster
  - `delete` - Delete a Kind cluster
  - `list` - List all Kind clusters

### Provider Management

- `crosslab provider` - Manage Crossplane providers
  - `install` - Install a specific provider
  - `install-all` - Install all required providers
  - `list` - List installed providers

## Kind Cluster Management

### Create a Cluster

```bash
crosslab cluster create --config examples/config/kind-config.yaml --name my-cluster
```

### Delete a Cluster

```bash
crosslab cluster delete --name my-cluster
```

### List Clusters

```bash
crosslab cluster list
```

### Example Kind Configuration

An example Kind cluster configuration is provided in `examples/config/kind-config.yaml`. This configuration creates a cluster with:
- 1 control plane node
- 2 worker nodes
- Port mappings for HTTP (80 → 8080) and HTTPS (443 → 8443)
- Custom pod and service subnets

## Crossplane Provider Management

### Provider Configuration

Providers are configured in a YAML file. The default location is `.crosslab/providers.yaml`, but you can specify a different file with the `--config` flag.

The configuration structure is:

```yaml
aws:
  family:
    name: string       # Provider name
    package: string    # Provider package
    version: string    # Provider version
  services:
    - name: string     # Service provider name
      package: string  # Service provider package
      version: string  # Service provider version

otherProviders:
  - name: string       # Provider name
    package: string    # Provider package
    version: string    # Provider version
```

An example configuration is available at `examples/config/crosslab-config.yaml`.

### Install a Specific Provider

```bash
# Install new provider
crosslab provider install \
  --name provider-aws \
  --package xpkg.upbound.io/upbound/provider-aws \
  --version v1.0.0

# Force reinstall existing provider
crosslab provider install \
  --name provider-aws \
  --package xpkg.upbound.io/upbound/provider-aws \
  --version v1.0.0 \
  --force
```

### Install All Required Providers

This command will install all providers defined in the configuration file:

```bash
# Install new providers using default config
crosslab provider install-all

# Install new providers using specific config
crosslab provider install-all --config path/to/providers.yaml

# Force reinstall all providers
crosslab provider install-all --force
```

### List Installed Providers

```bash
crosslab provider list
```

## Development

### Available Make Commands

#### Build and Run

- `make build` - Build the application using GoReleaser
- `make run` - Run the application
- `make dev` - Build and run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts

#### Code Quality

- `make fmt` - Format code
- `make lint` - Run linter

#### Module Management

- `make mod-download` - Download all dependencies
- `make mod-tidy` - Tidy the go.mod and go.sum files
- `make mod-verify` - Verify dependencies
- `make mod-update` - Update all dependencies
- `make mod-init` - Initialize a new module
- `make mod-graph` - Show module dependency graph

#### Release Management

- `make release` - Create a new release using GoReleaser
- `make install-goreleaser` - Install GoReleaser
- `make install-binary` - Install binary from GitHub releases

### Running the Project

Using Make:

```bash
make run
```

Or using Go directly:

```bash
go run main.go
```

### Supported Platforms

The CLI supports the following platforms:
- Linux (amd64, arm64)
- Windows (amd64, arm64)
- macOS (amd64, arm64)

### Release Process

The project uses GoReleaser for automated releases with the following workflow:
- Snapshot builds for main and feature branches
- Pre-release builds for release candidates
- Full releases for version tags

Release artifacts include:
- Binary builds for all supported platforms
- Checksums for verification
- Source code archives
- Changelog generation

## Project Structure

```
.
├── Makefile                # Build and development commands
├── README.md               # This file
├── cmd/                    # CLI commands
│   └── crosslab/           # Command implementations
├── pkg/                    # Core functionality
│   └── common/             # Shared utilities
├── examples/               # Example configurations
│   └── config/             # Configuration examples
├── .crosslab/              # Default configuration directory
├── go.mod                  # Go module definition
└── main.go                 # Application entry point
```
