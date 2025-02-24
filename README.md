# Crosslab

A Cobra-based CLI application with Kind cluster and Crossplane provider management capabilities.
This is a project to enable the creation of a local development environment for crossplane based work such as compostion creation, testing, and development.

# Installation

## Pre-built binaries

Pre-built binaries are available for download from the [releases page](https://github.com/opus2-platform/crosslab-cli/releases).

## Building from source

To build from source, you will need to have Go installed on your system. You can download the latest version of Go from the [Go website](https://go.dev/dl/).

Once Go is installed, you can clone the repository and build the project using the following commands:

```bash
git clone https://github.com/opus2-platform/crosslab-cli.git
cd crosslab-cli
make build
```


MacOS and Linux users can use the `crosslab` binary directly, while Windows users will need to use the `crosslab.exe` binary.

If you use the install script, you can preview what occurs during the install process:


To install, run:

```bash
curl -fsSL https://raw.githubusercontent.com/kanzifucius/crosslab-cli/main/install.sh | sh
```

## Getting Started

### Prerequisites

- Go 1.23.2 or later
- Make
- Docker (for Kind clusters)
- GoReleaser (optional, for releases)

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

#### CLI Development
- `make new-command` - Create a new CLI command (interactive)

### CLI Commands

Currently available commands:
- `crosslab` - Root command
- `crosslab version` - Show the CLI version
- `crosslab cluster` - Manage Kind clusters
  - `create` - Create a new Kind cluster
  - `delete` - Delete a Kind cluster
  - `list` - List all Kind clusters
- `crosslab provider` - Manage Crossplane providers
  - `install` - Install a specific provider
  - `install-all` - Install all required providers
  - `list` - List installed providers

### Kind Cluster Management

The CLI provides commands to manage Kind clusters:

#### Create a Cluster
```bash
crosslab cluster create --config examples/kind-config.yaml --name my-cluster
```

#### Delete a Cluster
```bash
crosslab cluster delete --name my-cluster
```

#### List Clusters
```bash
crosslab cluster list
```

### Crossplane Provider Management

The CLI provides commands to manage Crossplane providers:

#### Provider Configuration

Providers are configured in a YAML file (`config/providers.yaml`). The configuration structure is:

```yaml
aws:
  family:
    name: string       # Provider name
    package: string    # Provider package
    version: string    # Provider version
  services:
    - name: string    # Service provider name
      package: string # Service provider package
      version: string # Service provider version

otherProviders:
  - name: string      # Provider name
    package: string   # Provider package
    version: string   # Provider version
```

Example configuration:
```yaml
aws:
  family:
    name: "upbound-provider-aws"
    package: "xpkg.upbound.io/upbound/provider-family-aws"
    version: "v1"
  services:
    - name: "provider-aws-iam"
      package: "xpkg.upbound.io/upbound/provider-aws-iam"
      version: "v1"

otherProviders:
  - name: "provider-helm"
    package: "xpkg.upbound.io/upbound/provider-helm"
    version: "v0.20.4"
```

#### Install a Specific Provider
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

#### Install All Required Providers
This command will install all providers defined in the configuration file:
- AWS Provider (family)
- AWS Service Providers (as configured)
- Other Providers (as configured)

```bash
# Install new providers using default config
crosslab provider install-all

# Install new providers using specific config
crosslab provider install-all --config path/to/providers.yaml

# Force reinstall all providers
crosslab provider install-all --force

# Force reinstall using specific config
crosslab provider install-all --config path/to/providers.yaml --force
```

#### List Installed Providers
```bash
crosslab provider list
```

### Example Kind Configuration

An example Kind cluster configuration is provided in `examples/kind-config.yaml`. This configuration creates a cluster with:
- 1 control plane node
- 2 worker nodes
- Port mappings for HTTP (80) and HTTPS (443)
- Custom pod and service subnets

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

### Project Structure

```
.
├── Makefile
├── README.md
├── cmd/
│   └── crosslab/
│       ├── root.go      # Root command
│       ├── version.go   # Version command
│       ├── cluster.go   # Cluster management commands
│       └── provider.go  # Provider management commands
├── pkg/
│   └── common/
│       ├── kind.go           # Kind cluster utilities
│       ├── crossplane.go     # Crossplane utilities
│       ├── config.go         # Configuration loader
│       ├── config_types.go   # Configuration types
│       └── provider_types.go # Provider type
├── config/
│   └── providers.yaml   # Provider configuration
├── examples/
│   └── kind-config.yaml # Example Kind configuration
├── go.mod
└── main.go
```

### VS Code Integration

The project includes VS Code configurations for development:

### Running the Project

Using Make:

```bash
make run
```

Or using Go directly:

```bash
go run main.go
```
