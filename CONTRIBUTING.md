# Contributing to GopenAPI

Thank you for your interest in contributing to GopenAPI! We welcome contributions from the community.

## 🤝 How to Contribute

### Reporting Issues

- Check if the issue already exists in [GitHub Issues](https://github.com/shubhamku044/gopenapi/issues)
- Use the issue templates when available
- Provide clear reproduction steps
- Include your Go version, OS, and GopenAPI version

### Suggesting Features

- Open a [GitHub Discussion](https://github.com/shubhamku044/gopenapi/discussions) first
- Describe the use case and benefit
- Check if it aligns with the project goals

### Pull Requests

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Make** your changes following our guidelines
4. **Test** your changes thoroughly
5. **Commit** with clear messages
6. **Push** and create a Pull Request

## 🛠️ Development Setup

### Prerequisites

- Go 1.21+
- Git
- Make (optional but recommended)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gopenapi.git
cd gopenapi

# Install dependencies
go mod download

# Run tests to verify setup
make test
```

### Development Workflow

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build binary
make build

# Run all quality checks
make check
```

## 📋 Development Guidelines

### Code Style

- Follow standard Go conventions
- Run `go fmt` and `goimports`
- Use meaningful variable and function names
- Add comments for public APIs
- Keep functions small and focused

### Testing

- Write unit tests for new functionality
- Maintain or improve test coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate
- Mock external dependencies

### Commits

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or modifications
- `chore`: Build process or auxiliary tool changes

Examples:
```
feat(generator): add support for path parameters
fix(parser): handle empty operation IDs correctly
docs(readme): update installation instructions
```

### Code Organization

```
gopenapi/
├── cmd/gopenapi/          # CLI application
├── internal/              # Internal packages
│   ├── generator/         # Code generation logic
│   ├── models/           # Data models
│   └── parser/           # OpenAPI parsing
├── pkg/                  # Public packages
│   ├── templates/        # Code templates
│   └── utils/           # Utility functions
├── example/             # Example usage
└── .github/            # GitHub workflows
```

### Architecture Principles

- **Separation of Concerns**: Keep parsing, generation, and CLI logic separate
- **Testability**: Write testable code with clear interfaces
- **Maintainability**: Prefer clarity over cleverness
- **Performance**: Optimize for reasonable performance, not micro-optimizations
- **Compatibility**: Maintain backward compatibility when possible

## 🧪 Testing

### Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# With race detection
make test-race

# Specific package
go test ./internal/generator

# Specific test
go test -run TestGenerateCode ./internal/generator
```

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete workflows
4. **CLI Tests**: Test command-line interface

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    t.Run("descriptive test case name", func(t *testing.T) {
        // Arrange
        input := "test input"
        expected := "expected output"
        
        // Act
        result, err := FunctionName(input)
        
        // Assert
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if result != expected {
            t.Errorf("got %q, want %q", result, expected)
        }
    })
}
```

## 📦 Release Process

Releases are automated via GitHub Actions when tags are pushed:

```bash
# Create and push a new tag
git tag v1.2.3
git push origin v1.2.3
```

This triggers:
- Cross-platform binary builds
- GitHub release creation
- Archive generation
- Changelog updates

## 🏷️ Labels and Issues

We use these labels for organization:

- `bug`: Something isn't working
- `enhancement`: New feature or request
- `documentation`: Improvements or additions to docs
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention is needed
- `priority/high`: High priority
- `priority/low`: Low priority

## 📖 Documentation

- Update README.md for user-facing changes
- Update DOCUMENTATION.md for detailed technical docs
- Add inline comments for complex logic
- Update examples when adding features

## ❓ Questions?

- 💬 [GitHub Discussions](https://github.com/shubhamku044/gopenapi/discussions)
- 🐛 [GitHub Issues](https://github.com/shubhamku044/gopenapi/issues)
- 📧 Email: [your-email@example.com]

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for helping make GopenAPI better! 🙏 