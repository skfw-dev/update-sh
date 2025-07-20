# Security Policy

## Supported Versions

We provide security updates for the following versions of Update-SH:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

### Private Disclosure Process

We take security issues seriously and appreciate your efforts to responsibly disclose your findings. To report a security issue, please follow these steps:

1. **Do not** create a public GitHub issue
2. Email our security team at [security@skfw.dev](mailto:security@skfw.dev) with:
   - A detailed description of the vulnerability
   - Steps to reproduce the issue
   - Any proof-of-concept code
   - Your contact information
3. We will acknowledge receipt of your report within 48 hours
4. We will review the report and determine the impact and severity
5. We will provide a timeline for addressing the issue
6. Once resolved, we will release a security advisory and credit you (unless you prefer to remain anonymous)

### Public Disclosure Timeline

After a security vulnerability has been addressed, we will:
- Release a security patch for all supported versions
- Publish a security advisory on GitHub
- Credit the reporter (unless requested otherwise)
- Update the documentation to reflect any security-related changes

## Security Best Practices

### For Users
- Always use the latest stable version of Update-SH
- Run the tool with the minimum required privileges
- Review and understand the packages being updated
- Use the `--dry-run` flag to preview changes before applying them

### For Developers
- Follow secure coding practices
- Keep dependencies up to date
- Use the latest stable version of Go
- Run security scanners as part of the CI/CD pipeline
- Document security-related code clearly

## Security Updates

Security updates are released as patch versions (e.g., 1.0.0 → 1.0.1). We recommend always running the latest patch version of your installed major.minor version.

## Security Advisories

Security advisories are published on our [GitHub Security Advisories](https://github.com/skfw-dev/update-sh/security/advisories) page.
