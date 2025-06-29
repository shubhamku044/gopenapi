# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GitHub Actions workflow for automated releases
- Cross-platform binary builds (Linux, macOS, Windows)
- Comprehensive test suite with 76%+ coverage
- golangci-lint integration for code quality
- Integration tests for end-to-end validation

### Changed
- Updated README with installation instructions
- Improved project documentation

## [v0.1.0] - 2024-01-XX

### Added
- Initial release of GopenAPI
- OpenAPI 3.0 specification parsing
- Go code generation with Gin framework support
- Clean separation between generated and user code
- Safe regeneration without overwriting custom code
- Support for:
  - Data models from OpenAPI schemas
  - HTTP handlers with Gin routing
  - API interfaces for easy testing
  - Server boilerplate with graceful shutdown
  - Automatic README documentation generation
- Command-line interface with options:
  - `--spec`: OpenAPI specification file path
  - `--output`: Output directory
  - `--package`: Package name for generated code
- Project structure generation:
  - `main.go` - Application entry point
  - `go.mod` - Module definition
  - `handlers/` - User business logic
  - `generated/` - Generated code (safe to overwrite)
- Examples and documentation
- MIT License

### Features
- ✅ Clean code generation
- ✅ Safe regeneration
- ✅ Gin framework support
- ✅ Type safety from OpenAPI schemas
- ✅ Production-ready server setup
- ✅ Comprehensive documentation
- ✅ Zero runtime dependencies

[Unreleased]: https://github.com/shubhamku044/gopenapi/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/shubhamku044/gopenapi/releases/tag/v0.1.0 