#!/bin/bash

# Build the gopenapi tool
echo "Building gopenapi..."
go build -o bin/gopenapi ./cmd/gopenapi

# Install the binary
echo "Installing gopenapi..."
if [ -w "/usr/local/bin" ]; then
    # If we have write permissions to /usr/local/bin
    mv bin/gopenapi /usr/local/bin/gopenapi
    echo "Installation complete. You can now use 'gopenapi' from anywhere."
else
    # Use sudo if we don't have write permissions
    echo "Requires administrator privileges to install to /usr/local/bin"
    sudo mv bin/gopenapi /usr/local/bin/gopenapi
    echo "Installation complete. You can now use 'gopenapi' from anywhere."
fi