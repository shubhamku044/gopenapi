#!/bin/bash

# Build the gopenapi tool
echo "Building gopenapi..."
go build -o bin/gopenapi ./cmd/gopenapi

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Running code generation..."
    ./bin/gopenapi generate all --input=petstore.yaml --output=./gen
else
    echo "Build failed!"
fi
