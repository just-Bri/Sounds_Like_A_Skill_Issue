#!/bin/bash

# Game launcher script for SLASI
GAME_NAME="SLASI_arm_v1"
GAME_DIR="$(dirname "$(readlink -f "$0")")"
GAME_PATH="$GAME_DIR/$GAME_NAME"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print error message and exit
error_exit() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

# Function to print info message
info() {
    echo -e "${GREEN}[INFO] $1${NC}"
}

# Function to print warning message
warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

# Check if the game executable exists
if [ ! -f "$GAME_PATH" ]; then
    error_exit "Game executable not found at: $GAME_PATH"
fi

# Check if the file is executable, if not make it executable
if [ ! -x "$GAME_PATH" ]; then
    info "Setting executable permissions for the game..."
    chmod +x "$GAME_PATH" || error_exit "Failed to set executable permissions"
fi

# Check for required libraries (add more as needed)
for lib in "libasound.so" "libX11.so" "libXrandr.so" "libXcursor.so" "libXi.so"; do
    if ! ldconfig -p | grep -q "$lib"; then
        warning "Library $lib might be missing. Game may not run correctly."
    fi
done

# Set up environment variables for better compatibility
export SDL_AUDIODRIVER=alsa
export MESA_GL_VERSION_OVERRIDE=3.3

# Optional: Check for gamepad/controller
if [ -d "/dev/input" ]; then
    if ls /dev/input/js* 1> /dev/null 2>&1; then
        info "Gamepad detected"
    else
        warning "No gamepad detected. Game can still be played with keyboard."
    fi
fi

# Launch the game
info "Launching SLASI..."
cd "$GAME_DIR" # Change to game directory before launching
"$GAME_PATH"

# Store the exit code
exit_code=$?

# Check if game exited normally
if [ $exit_code -ne 0 ]; then
    error_exit "Game exited with error code: $exit_code"
else
    info "Game closed successfully"
fi

exit 0
