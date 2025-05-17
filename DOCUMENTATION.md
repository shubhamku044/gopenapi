# gopenapi Documentation

## Overview

gopenapi is a lightweight, opinionated code generator that reads a custom YAML/JSON specification file and generates idiomatic Go code including:

- Structs for models
- HTTP handler stubs with routing
- HTTP clients for consuming the API

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────────┐
│                 │     │                 │     │                     │
│  gopenapi.yaml  │────▶│  gopenapi CLI   │────▶│  Generated Go Code  │
│  specification  │     │                 │     │                     │
│                 │     │                 │     │                     │
└─────────────────┘     └─────────────────┘     └─────────────────────┘
                              │
                              │
                              ▼
                        ┌─────────────────┐
                        │                 │
                        │    Generator    │
                        │                 │
                        └─────────────────┘
                              │
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
    ┌─────────────────┐┌─────────────────┐┌─────────────────┐
    │                 ││                 ││                 │
    │  Model Generator││ Server Generator││ Client Generator│
    │                 ││                 ││                 │
    └─────────────────┘└─────────────────┘└─────────────────┘
```

## How It Works

1. **Input Specification**: You define your API in a custom YAML/JSON format (gopenapi format)
2. **Code Generation**: The CLI tool parses the specification and generates Go code
3. **Output**: Generated code includes models, server handlers, and client code

## gopenapi Format

The gopenapi format is a custom, simplified alternative to the OpenAPI specification. It's designed to be more intuitive and easier to work with for Go developers.

### Example

```yaml
gopenapi: 1.0.0
api:
  name: Petstore API
  version: 1.0.0
  description: A simple petstore API example

endpoints:
  - path: /pets
    methods:
      - type: GET
        name: listPets
        summary: List all pets
        response:
          status: 200
          description: A list of pets
          type: array
          items: Pet
      - type: POST
        name: createPet
        summary: Create a pet
        request:
          type: Pet
        response:
          status: 201
          description: Pet created
  - path: /pets/{petId}
    methods:
      - type: GET
        name: getPet
        summary: Get a pet by ID
        parameters:
          - name: petId
            in: path
            required: true
            type: integer
        response:
          status: 200
          description: A pet
          type: Pet

models:
  - name: Pet
    description: A pet in the store
    required:
      - id
      - name
    properties:
      - name: id
        type: integer
        format: int64
      - name: name
        type: string
      - name: tag
        type: string
      - name: status
        type: string
        enum:
          - available
          - pending
          - sold
```

## Generated Code

### Models

The generator creates Go structs for each model defined in the specification:

```go
// Pet A pet in the store
type Pet struct {
    Id int64 `json:"id,omitempty" yaml:"id,omitempty" validate:"required"`
    Name string `json:"name,omitempty" yaml:"name,omitempty" validate:"required"`
    Tag string `json:"tag,omitempty" yaml:"tag,omitempty"`
    Status string `json:"status,omitempty" yaml:"status,omitempty"`
}
```

### Client

The generator creates an HTTP client with methods for each endpoint:

```go
// Client is an HTTP client for the Petstore API
type Client struct {
    BaseURL    string
    HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL:    baseURL,
        HTTPClient: http.DefaultClient,
    }
}

// listPets List all pets
func (c *Client) listPets() ([]models.Pet, error) {
    // Implementation...
}

// createPet Create a pet
func (c *Client) createPet(req *models.Pet) error {
    // Implementation...
}

// getPet Get a pet by ID
func (c *Client) getPet(petId int) (*models.Pet, error) {
    // Implementation...
}
```

### Server

The generator creates HTTP handlers and routing for each endpoint:

```go
// Implementation will be available in future versions
```

## Command Line Usage

```bash
# Generate models
./gopenapi generate model --input=petstore.yaml --output=./gen

# Generate client code
./gopenapi generate client --input=petstore.yaml --output=./gen

# Generate server code
./gopenapi generate server --input=petstore.yaml --output=./gen

# Generate all components
./gopenapi generate all --input=petstore.yaml --output=./gen
```

## Internal Architecture

The generator is composed of several key components:

1. **Parser**: Reads and validates the YAML/JSON specification
2. **Model Generator**: Generates Go structs from model definitions
3. **Client Generator**: Generates HTTP client code from endpoint definitions
4. **Server Generator**: Generates HTTP handlers and routing from endpoint definitions

Each generator uses Go templates to produce the output code, making it easy to customize the generated code if needed.

## Compatibility

gopenapi also supports the standard OpenAPI 3.x format as a fallback, so you can use existing OpenAPI specifications with this tool.
