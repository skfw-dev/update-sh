#!/usr/bin/env bash

# set -e: Exit immediately if a command exits with a non-zero status.
# set -u: Treat unset variables as an error when substituting.
# set -o pipefail: The return value of a pipeline is the status of the last command
# to exit with a non-zero status, or zero if no command exited with a non-zero status.
set -euo pipefail

# --- Configuration ---
# Define variables for key paths and file names to make the script easier to modify.
readonly BINARY_NAME="update-sh"
readonly SOURCE_DIR="."
readonly BUILD_FLAGS="-ldflags='-s -w'"
readonly DEST_DIR="/usr/local/bin" # Use /usr/local/bin for user-installed binaries.

# --- Functions ---
# A function to handle cleanup tasks on script exit, regardless of success or failure.
cleanup() {
    echo "Cleaning up build artifacts..."
    rm -f "${BINARY_NAME}"
}

# Trap signals (e.g., EXIT, ERR) to ensure the cleanup function is always called.
# This prevents the built binary from being left behind in the current directory.
trap cleanup EXIT ERR

# --- Validation and Pre-flight Checks ---

echo "--- Starting build and installation process for ${BINARY_NAME} ---"

# Check for required dependencies. 'go' must be in the user's PATH.
if ! command -v go &> /dev/null; then
    echo "Error: 'go' command not found. Please ensure Go is installed and in your PATH." >&2
    exit 1
fi

# Ensure the build directory exists and contains a Go module.
if [[ ! -f "${SOURCE_DIR}/go.mod" ]]; then
    echo "Error: Go module file 'go.mod' not found in '${SOURCE_DIR}'." >&2
    exit 1
fi

# Check if the destination directory is writable by the current user.
if [[ ! -w "${DEST_DIR}" ]]; then
    echo "Error: Destination directory '${DEST_DIR}' is not writable. Please run with 'sudo'." >&2
    exit 1
fi

# --- Build Process ---
echo "Building Go binary..."
echo "Running: go build -o ${BINARY_NAME} ${BUILD_FLAGS} ${SOURCE_DIR}"
bash -c "go build -o \"${BINARY_NAME}\" ${BUILD_FLAGS} \"${SOURCE_DIR}\""

# Check the exit status of the 'bash -c' command.
# $? holds the exit code of the last executed command.
if [ $? -ne 0 ]; then
    echo "Error: The Go build command failed. Exiting." >&2
    exit 1
fi

# Validate the build. Check if the binary was successfully created.
if [[ ! -f "${BINARY_NAME}" ]]; then
    echo "Error: Go build failed. The binary '${BINARY_NAME}' was not created." >&2
    exit 1
fi

echo "Build successful! Binary created at ./${BINARY_NAME}"

# --- Installation Process ---
echo "Installing binary to ${DEST_DIR}..."

# Check if an existing version of the binary is present.
if [[ -f "${DEST_DIR}/${BINARY_NAME}" ]]; then
    echo "A previous version exists. Removing it..."
    # A more robust approach might be to back it up first.
    rm -f "${DEST_DIR}/${BINARY_NAME}" || { echo "Error: Failed to remove old binary." >&2; exit 1; }
fi

# Move the newly built binary to the destination directory.
mv "${BINARY_NAME}" "${DEST_DIR}/${BINARY_NAME}" || { echo "Error: Failed to move binary." >&2; exit 1; }

# Set executable permissions for all users.
chmod a+x "${DEST_DIR}/${BINARY_NAME}" || { echo "Error: Failed to set permissions." >&2; exit 1; }

echo "--- Installation of ${BINARY_NAME} completed successfully! ---"