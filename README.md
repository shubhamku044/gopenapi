# GopenAPI

[![CI](https://github.com/shubhamku044/gopenapi/actions/workflows/ci.yml/badge.svg)](https://github.com/shubhamku044/gopenapi/actions/workflows/ci.yml)
[![Release](https://github.com/shubhamku044/gopenapi/actions/workflows/release.yml/badge.svg)](https://github.com/shubhamku044/gopenapi/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/shubhamku044/gopenapi)](https://goreportcard.com/report/github.com/shubhamku044/gopenapi)
[![GoDoc](https://godoc.org/github.com/shubhamku044/gopenapi?status.svg)](https://godoc.org/github.com/shubhamku044/gopenapi)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**A modern OpenAPI code generator for Go with Gin framework support.**

GopenAPI generates production-ready Go servers from OpenAPI 3.0 specifications with clean separation between generated and user code, ensuring your customizations are never overwritten.

## ✨ Features

- 🚀 **Clean Code Generation** - Separate generated and user code
- 🔄 **Safe Regeneration** - Never overwrites your custom code
- 🎯 **Gin Framework** - Built-in support for Gin HTTP router
- 📝 **Type Safety** - Strong typing from OpenAPI schemas
- 🛡️ **Production Ready** - Graceful shutdown, middleware support
- 📚 **Auto Documentation** - Generates comprehensive README
- 🧪 **Well Tested** - 76%+ test coverage

## 📦 Installation

### Option 1: Go Install (Recommended)
```bash
go install github.com/shubhamku044/gopenapi/cmd/gopenapi@latest
```

### Option 2: Download Binary
Download the latest binary from [GitHub Releases](https://github.com/shubhamku044/gopenapi/releases):

```bash
# Linux/macOS
curl -L https://github.com/shubhamku044/gopenapi/releases/latest/download/gopenapi_linux_amd64.tar.gz | tar xz

# macOS (ARM)
curl -L https://github.com/shubhamku044/gopenapi/releases/latest/download/gopenapi_darwin_arm64.tar.gz | tar xz

# Windows
# Download gopenapi_windows_amd64.zip from releases page
```

### Option 3: Homebrew (Coming Soon)
```bash
brew install shubhamku044/tap/gopenapi
```

### Option 4: Build from Source
```bash
git clone https://github.com/shubhamku044/gopenapi.git
cd gopenapi
go build -o gopenapi ./cmd/gopenapi
```

## 🚀 Quick Start

### 1. Create an OpenAPI spec (`api.yaml`)
```yaml
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths:
  /users:
    get:
      operationId: listUsers
      summary: List users
      responses:
        '200':
          description: OK
    post:
      operationId: createUser  
      summary: Create user
      responses:
        '201':
          description: Created
  /users/{id}:
    get:
      operationId: getUser
      summary: Get user by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: OK
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
        name:
          type: string
        email:
          type: string
          format: email
```

### 2. Generate your API server
```bash
gopenapi --spec=api.yaml --output=myapi --package=myapi
```

### 3. Implement your business logic
Edit `myapi/handlers/api.go`:
```go
func (h *APIHandlers) ListUsers(c *gin.Context) {
    users := []models.User{
        {Id: "1", Name: "John Doe", Email: "john@example.com"},
        {Id: "2", Name: "Jane Smith", Email: "jane@example.com"},
    }
    c.JSON(http.StatusOK, users)
}

func (h *APIHandlers) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Save user to database...
    user.Id = "generated-id"
    
    c.JSON(http.StatusCreated, user)
}

func (h *APIHandlers) GetUser(c *gin.Context, id string) {
    // Fetch user from database...
    user := models.User{Id: id, Name: "John Doe", Email: "john@example.com"}
    c.JSON(http.StatusOK, user)
}
```

### 4. Run your server
```bash
cd myapi
go mod tidy
go run main.go
```

Your API server is now running on `http://localhost:8080`! 🎉

## 🏗️ Project Structure

GopenAPI creates a clean project structure that separates your code from generated code:

```
myapi/
├── main.go              # ✏️  Your application entry point
├── go.mod              # ✏️  Your module definition  
├── handlers/           # ✏️  YOUR BUSINESS LOGIC
│   └── api.go         #     Implement your handlers here
├── generated/          # 🤖 Generated code (safe to regenerate)
│   ├── api/
│   │   └── interfaces.go  # API interface definitions
│   ├── models/
│   │   └── models.go      # Data models from OpenAPI spec
│   └── server/
│       └── router.go      # HTTP server and routing
└── README.md           # 📚 Generated documentation
```

**✅ Safe files** (never overwritten): `main.go`, `go.mod`, `handlers/`
**🔄 Generated files** (safe to regenerate): Everything in `generated/`

## 📖 Usage Examples

### Generate with custom package name
```bash
gopenapi --spec=api.yaml --output=./my-service --package=myservice
```

### Update existing project (safe regeneration)
```bash
gopenapi --spec=updated-api.yaml --output=. --package=myapi
```

### Help and options
```bash
gopenapi --help
```

## 🆚 Comparison with Other Tools

| Feature | GopenAPI | oapi-codegen | go-swagger |
|---------|----------|--------------|------------|
| **Safe Regeneration** | ✅ | ❌ | ❌ |
| **Gin Support** | ✅ | ✅ | ❌ |
| **Clean Separation** | ✅ | ❌ | ❌ |
| **Auto Documentation** | ✅ | ❌ | ✅ |
| **Type Safety** | ✅ | ✅ | ✅ |
| **Production Ready** | ✅ | ⚠️ | ✅ |

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup
```bash
git clone https://github.com/shubhamku044/gopenapi.git
cd gopenapi
go mod download

# Run tests
make test

# Run linter  
make lint

# Build binary
make build
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- Built with [Gin](https://github.com/gin-gonic/gin) web framework
- OpenAPI specification support

## 📞 Support

- 📚 [Documentation](https://github.com/shubhamku044/gopenapi/blob/main/DOCUMENTATION.md)
- 🐛 [Report Issues](https://github.com/shubhamku044/gopenapi/issues)
- 💬 [Discussions](https://github.com/shubhamku044/gopenapi/discussions)
- ⭐ Star this repo if you find it useful!

---

**Made with ❤️ by [Shubham Kumar](https://github.com/shubhamku044)**