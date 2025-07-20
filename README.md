# Update-SH: Cross-Platform System Update Manager

[![Go Report Card](https://goreportcard.com/badge/github.com/skfw-dev/update-sh)](https://goreportcard.com/report/github.com/skfw-dev/update-sh)
[![GitHub release](https://img.shields.io/github/v/release/skfw-dev/update-sh)](https://github.com/skfw-dev/update-sh/releases)
[![License: Apache 2.0](https://img.shields.io/github/license/skfw-dev/update-sh)](https://github.com/skfw-dev/update-sh/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/skfw-dev/update-sh.svg)](https://pkg.go.dev/github.com/skfw-dev/update-sh)
[![Build Status](https://github.com/skfw-dev/update-sh/actions/workflows/build.yml/badge.svg)](https://github.com/skfw-dev/update-sh/actions)
[![Codecov](https://codecov.io/gh/skfw-dev/update-sh/branch/main/graph/badge.svg)](https://codecov.io/gh/skfw-dev/update-sh)

A powerful command-line tool that simplifies system maintenance by managing updates across multiple package managers on both Linux and Windows platforms.

## 🚀 Features

- [x] **Cross-Platform Support**: Works on Linux and Windows
- [x] **Package Manager Integration**:
  - **Linux**: APT, DNF, Pacman, Zypper, Snap, Flatpak
  - **Windows**: WinGet, Chocolatey, Scoop
- [x] **Automatic Privilege Escalation**: Automatically requests admin/root privileges when needed
- [x] **Dry Run Mode**: Preview changes before applying them
- [x] **Verbose Output**: Detailed logging for troubleshooting
- [x] **Configurable**: Customize behavior via config file or command-line flags
- [x] **Shell Integration**: Optional updates for Zsh and PowerShell

## 📦 Installation

### Prerequisites

- **Linux/Windows**: Go 1.24 or later
- **Build Tools**: Git, Make (optional)

### Quick Install (Linux/macOS)

```bash
# Install with one command (requires sudo)
curl -sSL https://raw.githubusercontent.com/skfw-dev/update-sh/main/scripts/install.sh | bash
```

### Manual Installation

#### Linux

1. **Download the latest release**:
   ```bash
   # For amd64 systems
   curl -L -o update-sh https://github.com/skfw-dev/update-sh/releases/latest/download/update-sh-linux-amd64
   
   # Make it executable
   chmod +x update-sh
   
   # Move to a directory in your PATH
   sudo mv update-sh /usr/local/bin/
   ```

2. **Verify installation**:
   ```bash
   update-sh --version
   ```

#### Windows

1. **Download the latest release** from the [Releases](https://github.com/skfw-dev/update-sh/releases) page
2. **Add to PATH** or run from the downloaded location
3. **Run in PowerShell**:
   ```powershell
   .\update-sh.exe --help
   ```

### Using Go Install

```bash
go install github.com/skfw-dev/update-sh@latest
```

## 🛠️ Building from Source

### Prerequisites
- Go 1.24 or later
- Git
- Make (optional for development)

### Build Steps

#### Linux/macOS
```bash
git clone https://github.com/skfw-dev/update-sh.git
cd update-sh
make build  # or: go build -o bin/update-sh .
```

#### Windows
```powershell
git clone https://github.com/skfw-dev/update-sh.git
cd update-sh
go build -o bin/update-sh.exe .
```

## 🚀 Quick Start

### Basic Update
```bash
# Linux (will prompt for sudo)
sudo update-sh

# Windows (will request admin privileges)
update-sh
```

### Common Options
```bash
# Dry run (show what would be updated)
update-sh --dry-run

# Enable verbose output
update-sh -v

# Update Zsh and PowerShell components
update-sh --zsh-update --pwsh-update

# Specify a custom config file
update-sh --config ~/.config/update-sh/config.yaml
```

## ⚙️ Configuration

Create a config file at `~/.update-sh.yaml` (Linux/macOS) or `%USERPROFFILE%\.update-sh.yaml` (Windows):

```yaml
# Example config file
verbose: true
dry-run: false
log-file: /var/log/update-sh.log
zsh-update: true
pwsh-update: true
```

## 📚 Documentation

For detailed documentation, please visit our [documentation website](https://skfw-dev.github.io/update-sh/).

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on how to submit pull requests, report issues, or suggest new features.

## 🔒 Security

Please see our [Security Policy](SECURITY.md) for information about reporting security vulnerabilities.

## 📄 License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.

## 📝 Changelog

See [CHANGELOG.md](CHANGELOG.md) for a history of changes to this project.

## 📬 Contact

- **Issues**: [GitHub Issues](https://github.com/skfw-dev/update-sh/issues)
- **Discussions**: [GitHub Discussions](https://github.com/skfw-dev/update-sh/discussions)
- **Email**: [contact@skfw.dev](mailto:contact@skfw.dev)

## 🙏 Acknowledgments

- Thanks to all [contributors](https://github.com/skfw-dev/update-sh/graphs/contributors) who have helped improve this project.
- Inspired by various system update tools and package managers.