#!/bin/bash

# Configuration Switcher Launch Script

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="$SCRIPT_DIR/switcher"

if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: switcher binary not found at $BINARY_PATH"
    echo "Please run 'go build -o switcher' first"
    exit 1
fi

echo "Starting Configuration Switcher..."
echo "Use arrow keys to navigate, Enter to select, Esc to go back, Ctrl+C to quit"
echo ""

# Check if we're in a terminal
if [ -t 0 ] && [ -t 1 ]; then
    exec "$BINARY_PATH"
else
    echo "Error: This application requires a terminal with TTY support"
    echo "Please run this from a proper terminal session"
    exit 1
fi