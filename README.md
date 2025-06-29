# gopenapi

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**gopenapi** is a lightweight, opinionated OpenAPI code generator built specifically for Go developers. It reads an OpenAPI 3.x specification file (YAML/JSON) and generates clean, idiomatic Go code.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/shubhamku044/gopenapi.git
cd gopenapi

# Install gopenapi globally
make install

# Generate code from an OpenAPI spec
gopenapi --spec=api.yaml --output=./generated --package=myapi
```

## 🎯 What is gopenapi?

gopenapi generates production-ready Go code from OpenAPI 3.x specifications, including:

- **🏗️ Data Models** - Clean structs for all schemas with proper Go types
- **🔌 HTTP Handlers** - Handler stubs with proper routing and middleware support
- **📡 API Interfaces** - Well-defined interfaces for easy testing and mocking
- **🚀 Server Boilerplate** - Complete HTTP server setup with Gin framework
- **📚 Documentation** - Generated README with usage examples

All with **zero runtime dependencies**, **no code bloat**, and **100% Go-native** tooling.

## 💡 Why gopenapi?

Existing solutions like Swagger Codegen or OpenAPI Generator:
- ❌ Are heavyweight, written in Java, require complex configurations
- ❌ Generate bloated, non-idiomatic Go code
- ❌ Are difficult to customize and contribute to
- ❌ Have steep learning curves

Go developers want something:
- 🐹 **Simple** - One binary, zero dependencies
- 🧬 **Go-native** - Written in Go, for Go developers
- 🧰 **Easy to integrate** - Fits naturally into Go workflows
- 🧠 **Customizable** - Easy to understand and modify
- ⚡ **Fast** - Lightweight and performant

That's where **gopenapi** comes in.

## 📦 Installation

### Using Make (Recommended)

```bash
# Clone and install
git clone https://github.com/shubhamku044/gopenapi.git
cd gopenapi
make install
```

### Manual Installation

```bash
# Build from source
go build -o gopenapi ./cmd/gopenapi
sudo mv gopenapi /usr/local/bin/
```

### Verify Installation

```bash
gopenapi --help
```

## 🛠️ Usage

### Basic Usage

```bash
# Generate code with default settings
gopenapi --spec=petstore.yaml --output=./generated

# Generate with custom package name
gopenapi --spec=api.yaml --output=./myapi --package=petstore

# Generate from JSON specification
gopenapi --spec=openapi.json --output=./client --package=apiclient
```

### Command Line Options

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `--spec` | Path to OpenAPI specification file (YAML/JSON) | *required* | `api.yaml` |
| `--output` | Output directory for generated code | `./generated` | `./myapi` |
| `--package` | Package name for generated code | auto-detected | `petstore` |

### Example Workflow

```bash
# 1. Start with an OpenAPI spec file
cat > petstore.yaml << 'EOF'
openapi: 3.0.0
info:
  title: Pet Store API
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      responses:
        '200':
          description: List of pets
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Pet'
components:
  schemas:
    Pet:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        status:
          type: string
          enum: [available, pending, sold]
EOF

# 2. Generate Go code
gopenapi --spec=petstore.yaml --output=./petstore-api

# 3. Your generated code structure:
# petstore-api/
# ├── api/
# │   └── api.go           # API interfaces
# ├── models/
# │   └── models.go        # Data models (Pet struct)
# ├── server/
# │   └── server.go        # HTTP server setup
# └── README.md            # Usage documentation
```

## 🔧 Development

### Using Make (Recommended)

```bash
# Show all available commands
make help

# Set up development environment
make dev-setup

# Build the project
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Generate code from sample API
make generate-sample

# Run all quality checks
make check
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `build` | Build the gopenapi binary |
| `install` | Build and install globally |
| `test` | Run all tests |
| `test-coverage` | Run tests with coverage report |
| `fmt` | Format code with go fmt and goimports |
| `lint` | Run golangci-lint |
| `clean` | Clean build artifacts |
| `generate-sample` | Generate code from sample-api.yaml |
| `release` | Build for all platforms |
| `dev-setup` | Install development tools |

### Manual Development

```bash
# Install dependencies
go mod download

# Build
go build -o bin/gopenapi ./cmd/gopenapi

# Run tests
go test ./...

# Generate sample
./bin/gopenapi --spec=sample-api.yaml --output=./generated
```

## 📁 Generated Code Structure

When you run gopenapi, it creates a well-organized directory structure:

```
generated/
├── api/
│   └── api.go           # API interfaces and handler definitions
├── models/
│   └── models.go        # Data models and structs
├── server/
│   └── server.go        # HTTP server setup with Gin
└── README.md            # Usage documentation and examples
```

### Example Generated Code

For a `Pet` schema, gopenapi generates:

```go
// models/models.go
type Pet struct {
    ID     int64  `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
}

// api/api.go
type PetAPI interface {
    ListPets(c *gin.Context)
    CreatePet(c *gin.Context)
    GetPet(c *gin.Context)
}

// server/server.go
func NewServer() *gin.Engine {
    r := gin.Default()
    // Routes and middleware setup
    return r
}
```

## 🎨 Sample API

The repository includes a `sample-api.yaml` file demonstrating a complete API specification. Try it out:

```bash
# Generate code from the sample
make generate-sample

# Or manually:
gopenapi --spec=sample-api.yaml --output=./sample-generated
```

## 🏗️ Build & Release

### Local Build

```bash
# Build for current platform
make build

# Build with debug symbols
make build-debug

# Cross-compile for all platforms
make cross-compile
```

### Release

```bash
# Create release builds for all platforms
make release

# This creates:
# bin/releases/gopenapi-linux-amd64.tar.gz
# bin/releases/gopenapi-darwin-amd64.tar.gz
# bin/releases/gopenapi-windows-amd64.exe.tar.gz
# ... and more
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run benchmarks
make benchmark
```

## 🔍 Code Quality

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run linter
make lint

# Run all quality checks
make check
```

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Make** your changes
4. **Run** quality checks: `make check`
5. **Test** your changes: `make test`
6. **Commit** your changes: `git commit -m 'Add amazing feature'`
7. **Push** to the branch: `git push origin feature/amazing-feature`
8. **Open** a Pull Request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gopenapi.git
cd gopenapi

# Set up development environment
make dev-setup

# Make your changes and test
make check test
```

## 📋 Requirements

- **Go 1.22+**
- **Git** (for version information)

## 🆚 Comparison with Other Tools

| Feature | gopenapi | oapi-codegen | Swagger Codegen |
|---------|----------|--------------|-----------------|
| **Language** | Go | Go | Java |
| **Binary Size** | ~10MB | ~15MB | ~50MB + JVM |
| **Dependencies** | Zero runtime | Minimal | Heavy (Java ecosystem) |
| **Go Idioms** | ✅ Native | ✅ Good | ❌ Poor |
| **Customization** | ✅ Easy | 🟡 Moderate | ❌ Complex |
| **Performance** | ⚡ Fast | ⚡ Fast | 🐌 Slow |
| **Learning Curve** | 🟢 Gentle | 🟡 Moderate | 🔴 Steep |

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the simplicity of Go tooling
- Built for the Go community, by Go developers
- Thanks to all contributors who help make this tool better

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/shubhamku044/gopenapi/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/shubhamku044/gopenapi/discussions)
- 📖 **Documentation**: This README and generated docs

---

**Made with ❤️ for the Go community**