# Contributing to Update-SH

Thank you for your interest in contributing to Update-SH! We appreciate your time and effort in helping improve this project. This document outlines the process for contributing to the project and what to expect from the maintainers.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Pull Requests](#pull-requests)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Code Review Process](#code-review-process)
- [Community](#community)
- [Support](#support)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report any unacceptable behavior to the project maintainers.

## Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally
   ```bash
   git clone https://github.com/skfw-dev/update-sh.git
   cd update-sh
   ```
3. **Set up** the development environment
   - Install Go 1.20 or later
   - Install project dependencies
   - Build the project
     ```bash
     go build -o bin/update-sh .
     ```
4. **Create a branch** for your changes
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b bugfix/issue-number-description
   ```

## How to Contribute

### Reporting Bugs

- **Check existing issues** to see if the bug has already been reported
- If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/skfw-dev/update-sh/issues/new/choose)
- Be sure to include:
  - A clear, descriptive title
  - Steps to reproduce the issue
  - Expected vs. actual behavior
  - Screenshots if applicable
  - Your operating system and version
  - Any relevant logs or error messages

### Suggesting Enhancements

- Use the "Feature Request" issue template
- Clearly describe the enhancement and why it would be useful
- Include any relevant use cases or examples
- If possible, suggest a proposed implementation approach

### Your First Code Contribution

Looking for a good first issue? Check out issues labeled with `good first issue` to get started!

### Pull Requests

1. **Keep it small** - Each PR should address a single issue or add a single feature
2. **Update documentation** - Ensure the README and other relevant docs are updated
3. **Add tests** - New features and bug fixes should include tests
4. **Follow the coding style** - See [Code Style](#code-style) below
5. **Update CHANGELOG.md** - Document your changes in the "Unreleased" section

## Development Workflow

1. Create a feature branch from `main`
2. Make your changes
3. Run tests and linters
   ```bash
   go test ./...
   golangci-lint run
   ```
4. Commit your changes following the [commit message guidelines](#commit-message-guidelines)
5. Push to your fork and open a pull request

## Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Keep lines under 120 characters
- Document all exported functions and types
- Write clear, concise comments that explain the "why" not the "what"

## Testing

### Running Tests

Run the full test suite with:
```bash
go test -v ./...
```

### Test Coverage

To generate a test coverage report:
```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

### Writing Tests

1. **Unit Tests**:
   - Place test files in the same directory as the code they test
   - Use `_test.go` suffix for test files
   - Follow the pattern `func TestXxx(t *testing.T)`
   - Use table-driven tests where appropriate

2. **Integration Tests**:
   - Place in `tests/integration` directory
   - Tag with `// +build integration`
   - Run with: `go test -tags=integration ./...`

3. **Platform-Specific Tests**:
   - Use build tags (`// +build windows`, `// +build linux`)
   - Test platform-specific functionality separately

### Mocking

Use interfaces to mock dependencies in tests. Example:
```go
package mock

type PackageManager interface {
    Update() error
    List() ([]string, error)
}
```

### Test Helpers

Create test helpers in `testutils` package to reduce code duplication.

## Community

### Getting Help

- Check the [documentation](https://github.com/skfw-dev/update-sh/docs)
- Search [existing issues](https://github.com/skfw-dev/update-sh/issues)
- Join our [Discord/Slack community]() (coming soon!)

### Discussion Forums

- [GitHub Discussions](https://github.com/skfw-dev/update-sh/discussions) - For questions and discussions
- [GitHub Issues](https://github.com/skfw-dev/update-sh/issues) - For bug reports and feature requests

### Events

- Weekly community calls (TBD)
- Hackathons and sprints (announced on GitHub Discussions)

## Support

### Reporting Security Issues

Please report security issues to security@example.com. We'll respond promptly.

### Commercial Support

For commercial support options, please contact support@example.com.

### Sponsorship

Become a sponsor to support the project's development:
[GitHub Sponsors](https://github.com/sponsors/skfw-dev)

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types**:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools and libraries

**Example**:
```
feat(apt): add support for Ubuntu 22.04

Add specific handling for Ubuntu 22.04 package repositories

Fixes #123
```

## Code Review Process

1. A maintainer will review your PR and may suggest changes
2. Address any feedback and push updates to your branch
3. Once approved, a maintainer will squash and merge your PR
4. Your changes will be included in the next release

## Thank You!

Your contributions to open source, large or small, make great projects like this possible. Thank you for being part of our community! 🎉