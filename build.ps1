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

# --- Determine OS and set build variables ---
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

# --- Set environment variables and output file ---
$env:GOOS = $goos
$env:GOARCH = $goarch
$binaryName = "update-sh$outputExtension"
$outputFile = "bin/$binaryName"

# --- Build the Go application ---
Write-Host "Building Go application..."
go build -o $outputFile -ldflags "-s -w" .

if (-not (Test-Path -Path $outputFile)) {
    Write-Error "Build failed: The output file was not created."
    exit 1
}

Write-Host "Build successful! Binary created at $outputFile"

# --- Install the binary ---
Write-Host "Starting installation phase..."

# Check if an existing binary is in the user's PATH
$existingBinaryPath = Get-Command -Name "update-sh$outputExtension" -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source

if ($existingBinaryPath) {
    Write-Host "Existing binary found at: $existingBinaryPath"
    $installDir = Split-Path -Path $existingBinaryPath -Parent
} else {
    Write-Host "No existing binary found in PATH. Creating new install directory."
    $installDir = Join-Path -Path $env:USERPROFILE -ChildPath ".bin"
    
    # Create the directory if it doesn't exist
    if (-not (Test-Path -Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
        Write-Host "Created new directory: $installDir"
    }
}

$destinationPath = Join-Path -Path $installDir -ChildPath $binaryName

# Move the new binary to the destination
Move-Item -Path $outputFile -Destination $destinationPath -Force

Write-Host "Installation of '$binaryName' completed successfully to: $destinationPath"

# --- Add .bin to PATH (for Windows) ---
if ($IsWindows) {
    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if ($userPath -notlike "*$installDir*") {
        Write-Host "Adding $installDir to the user's PATH environment variable."
        $newPath = "$userPath;$installDir"
        [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        Write-Host "Please open a new PowerShell terminal for the changes to take effect."
    }
}
