# gopenapi

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**gopenapi** is a lightweight, opinionated OpenAPI code generator built specifically for Go developers. It reads an OpenAPI 3.x specification file (YAML/JSON) and generates clean, idiomatic Go code directly in your project.

## 🚀 Quick Start

```bash
# 1. Install gopenapi
git clone https://github.com/shubhamku044/gopenapi.git
cd gopenapi
make install

# 2. Create your Go project
mkdir my-api && cd my-api
go mod init github.com/myuser/my-api

# 3. Create your OpenAPI specification
cat > api.yaml << 'EOF'
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths:
  /health:
    get:
      operationId: health
      summary: Health check
      responses:
        '200':
          description: OK
EOF

# 4. Generate your API code
gopenapi --spec=api.yaml

# 5. Build and run
go mod tidy
go run main.go
```

## 🎯 What is gopenapi?

gopenapi generates production-ready Go code from OpenAPI 3.x specifications, including:

- **🏗️ Data Models** - Clean structs for all schemas with proper Go types
- **🔌 HTTP Handlers** - Handler templates for easy implementation
- **📡 API Interfaces** - Well-defined interfaces for easy testing and mocking
- **🚀 Server Boilerplate** - Complete HTTP server with graceful shutdown
- **📚 Documentation** - Generated README with usage examples

All with **zero runtime dependencies**, **no code bloat**, and **100% Go-native** tooling.

## 💡 Why gopenapi?

Existing solutions like Swagger Codegen or OpenAPI Generator:
- ❌ Are heavyweight, written in Java, require complex configurations
- ❌ Generate bloated, non-idiomatic Go code
- ❌ Are difficult to customize and contribute to
- ❌ Overwrite your custom implementations

Go developers want something:
- 🐹 **Simple** - One binary, zero dependencies
- 🧬 **Go-native** - Written in Go, for Go developers
- 🛡️ **Safe** - Never overwrites your custom code
- 🧠 **Intuitive** - Works with standard Go project structure
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

### The Natural Go Workflow

gopenapi follows the standard Go project workflow:

```bash
# 1. Create your project (standard Go way)
mkdir petstore-api && cd petstore-api
go mod init github.com/myuser/petstore-api

# 2. Create your OpenAPI spec in the project
# (create petstore.yaml in your project directory)

# 3. Generate code directly in your project
gopenapi --spec=petstore.yaml

# 4. Your project now has this structure:
# petstore-api/
# ├── go.mod                    # Your Go module (preserved)
# ├── main.go                   # Generated main (only if new)
# ├── petstore.yaml            # Your OpenAPI spec
# ├── handlers/                # YOUR implementation code
# │   └── api.go              # Implement your business logic here
# ├── generated/              # Generated code (safe to regenerate)
# │   ├── api/
# │   │   └── interfaces.go   # API interface definitions
# │   ├── models/
# │   │   └── models.go       # Data models
# │   └── server/
# │       └── router.go       # HTTP server & routing
# └── README.md               # Generated documentation
```

### Command Line Options

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `--spec` | Path to OpenAPI specification file (YAML/JSON) | *required* | `api.yaml` |
| `--output` | Output directory | `.` (current directory) | `./myproject` |
| `--package` | Package name for generated code | auto-detected from go.mod | `petstore` |

### Safe Regeneration

The magic of gopenapi is **safe regeneration**:

```bash
# Edit your OpenAPI spec
vim api.yaml

# Regenerate safely - your custom code is preserved!
gopenapi --spec=api.yaml

# Your implementations in handlers/ are NEVER touched
# Only generated/ directory is updated
```

### Example Workflow

```bash
# 1. Create project and spec
mkdir petstore && cd petstore
go mod init github.com/myuser/petstore

cat > api.yaml << 'EOF'
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

# 2. Generate your API
gopenapi --spec=api.yaml

# 3. Implement your business logic
cat > handlers/api.go << 'EOF'
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/myuser/petstore/generated/models"
)

type APIHandlers struct{}

func NewAPIHandlers() *APIHandlers {
    return &APIHandlers{}
}

func (h *APIHandlers) ListPets(c *gin.Context) {
    pets := []models.Pet{
        {ID: 1, Name: "Buddy", Status: "available"},
        {ID: 2, Name: "Max", Status: "pending"},
    }
    c.JSON(http.StatusOK, pets)
}
EOF

# 4. Build and run
go mod tidy
go run main.go

# 5. Update your API spec and regenerate safely
# Edit api.yaml, then:
gopenapi --spec=api.yaml  # Your handlers/api.go is preserved!
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

# Demonstrate the complete workflow
make demo-workflow

# Run the example API server
make run-example

# Run all quality checks
make check
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `build` | Build the gopenapi binary |
| `install` | Build and install globally |
| `demo-workflow` | Demonstrate complete user workflow |
| `run-example` | Build and run example API server |
| `test` | Run all tests |
| `test-coverage` | Run tests with coverage report |
| `fmt` | Format code with go fmt and goimports |
| `lint` | Run golangci-lint |
| `clean` | Clean build artifacts |
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

# Try the example
cd example
../bin/gopenapi --spec=sample-api.yaml
```

## 📁 Generated Code Structure

When you run gopenapi in your project, it creates this structure:

```
your-project/
├── go.mod                    # YOUR Go module (never touched)
├── main.go                   # Generated only if new project
├── api.yaml                 # YOUR OpenAPI specification
├── handlers/                # YOUR implementation code
│   └── api.go              # Implement your business logic here
├── generated/              # Generated code (safe to overwrite)
│   ├── api/
│   │   └── interfaces.go   # API interface definitions
│   ├── models/
│   │   └── models.go       # Data models and structs
│   └── server/
│       └── router.go       # HTTP server setup with Gin
└── README.md               # Generated usage documentation
```

### Separation of Concerns

- **`handlers/`** - Your business logic (never overwritten)
- **`generated/`** - Generated boilerplate (safe to regenerate)
- **`main.go`** - Application entry point (generated once)
- **Root files** - Your project files (go.mod, specs, etc.)

### Example Generated Code

For a `Pet` schema, gopenapi generates:

```go
// generated/models/models.go
type Pet struct {
    ID     *int64  `json:"id,omitempty"`
    Name   *string `json:"name,omitempty"`
    Status *string `json:"status,omitempty"`
}

// generated/api/interfaces.go
type API interface {
    ListPets(c *gin.Context)
    CreatePet(c *gin.Context)
    GetPet(c *gin.Context)
}

// generated/server/router.go
func NewServer(api api.API) *gin.Engine {
    r := gin.Default()
    r.GET("/pets", api.ListPets)
    r.POST("/pets", api.CreatePet)
    r.GET("/pets/:id", api.GetPet)
    return r
}

// handlers/api.go (your implementation)
func (h *APIHandlers) ListPets(c *gin.Context) {
    // YOUR business logic here
    c.JSON(200, []models.Pet{{ID: &[]int64{1}[0], Name: &[]string{"Buddy"}[0]}})
}
```

## 🎨 Example Project

The repository includes a complete example in the `example/` directory:

```bash
# Try the example
make run-example

# Or manually:
cd example
gopenapi --spec=sample-api.yaml
go mod tidy
go run main.go
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
| **Safe Regen** | ✅ Yes | ❌ No | ❌ No |
| **Project Integration** | ✅ Native | 🟡 Moderate | ❌ Complex |
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