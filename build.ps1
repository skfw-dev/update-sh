#!pwsh

# Build script for update-sh
# This script automates the build process for the update-sh project,
# including running tests and generating documentation.
# Usage: .\build.ps1

param (
    [string]$goos = "",
    [string]$goarch = "",
    [string]$outputExtension = ""
)

if ($IsLinux) {
    $goos = "linux"
    $goarch = "amd64"
    $outputExtension = ""

  } elseif ($IsWindows) {
    $goos = "windows"
    $goarch = "amd64"
    $outputExtension = ".exe"

  } elseif ($IsMacOS) {
    $goos = "darwin"
    $goarch = "amd64"
    $outputExtension = ""

  } else {
    Write-Error "Unsupported OS"
    exit 1
}

# Set the GOOS and GOARCH environment variables
$env:GOOS = $goos
$env:GOARCH = $goarch

# Set the output file
$outputFile = "bin/update-sh-$goos-$goarch"
if ($outputExtension) {
    $outputFile = "bin/update-sh-$goos-$goarch$outputExtension"
}

# Build the Go application
go build -o $outputFile -ldflags "-s -w" .
