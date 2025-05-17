# gopenapi

gopenapi is a command-line tool for generating API code boilerplate from YAML/JSON files

## What is gopenapi?

gopenapi is a lightweight, opinionated OpenAPI code generator built specifically for Go developers. It reads an OpenAPI 3.x spec file (YAML/JSON) and generates idiomatic Go code including:

- Structs for schemas (models)
- HTTP handler stubs with routing
- Request/response models with proper types
- Optionally, Go HTTP clients for consuming the API

All with no code bloat, no foreign tooling (e.g., Java), and zero runtime dependencies.

⸻

💡 Why does this project exist?

Existing solutions like Swagger Codegen or OpenAPI Generator:
- Are heavyweight, written in Java, and require complex configs
- Generate bloated, non-idiomatic Go code
- Are hard to customize and contribute to

Go developers want something:
- 🐹 Simple
- 🧬 Go-native
- 🧰 Easy to integrate into their workflow
- 🧠 Customizable if needed

That’s where gopenapi comes in.

⸻

## How Developers Will Use It

Simple CLI usage:
```bash
# Generate Go structs from OpenAPI spec
gopenapi generate model --input=openapi.yaml --output=./gen/models

# Generate handler + routing stubs
gopenapi generate server --input=openapi.yaml --output=./gen

# All in one
gopenapi generate all --input=openapi.yaml --output=./gen

```