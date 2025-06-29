package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shubhamku044/gopenapi/internal/generator"
	"github.com/shubhamku044/gopenapi/internal/parser"
)

func main() {
	specFile := flag.String("spec", "", "Path to OpenAPI specification file (YAML or JSON)")
	outputDir := flag.String("output", ".", "Output directory for generated code (defaults to current directory)")
	packageName := flag.String("package", "", "Package name for generated code (auto-detected from go.mod if not provided)")
	flag.Parse()

	if *specFile == "" {
		log.Fatal("Please provide an OpenAPI specification file with --spec")
	}

	// Check if spec file exists
	if _, err := os.Stat(*specFile); os.IsNotExist(err) {
		log.Fatalf("OpenAPI specification file '%s' not found", *specFile)
	}

	spec, err := parser.ParseSpecFile(*specFile)
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI specification: %v", err)
	}

	// Always try to read module name from go.mod first
	moduleName := ""
	goModPath := filepath.Join(*outputDir, "go.mod")
	if content, err := os.ReadFile(goModPath); err == nil {
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "module ") {
				moduleName = strings.TrimSpace(line[7:])
				break
			}
		}
	}

	// Determine package name
	pkg := *packageName
	if pkg == "" {
		if moduleName != "" {
			// Extract package name from module (last part after /)
			pkg = filepath.Base(moduleName)
		} else {
			// Fallback to directory name if no go.mod found
			absPath, err := filepath.Abs(*outputDir)
			if err == nil {
				pkg = filepath.Base(absPath)
			}
			if pkg == "" || pkg == "." {
				pkg = "api"
			}
		}
	}

	// Use package name as module name fallback if no go.mod found
	if moduleName == "" {
		moduleName = pkg
	}

	config := generator.Config{
		OutputDir:   *outputDir,
		PackageName: pkg,
		ModuleName:  moduleName,
	}

	err = generator.GenerateCode(spec, config)
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	fmt.Printf("✅ Code successfully generated in %s\n", *outputDir)
	fmt.Printf("📦 Module: %s\n", moduleName)
	fmt.Printf("📁 Package: %s\n", pkg)

	// Check if dependencies need to be installed
	goModExists := false
	if _, err := os.Stat(filepath.Join(*outputDir, "go.mod")); err == nil {
		goModExists = true
	}

	fmt.Println("\n🚀 Next steps:")
	if goModExists {
		fmt.Println("   go mod tidy                    # Install/update dependencies")
	} else {
		fmt.Println("   go mod init <your-module-name> # Initialize Go module first")
		fmt.Println("   go mod tidy                    # Install dependencies")
	}
	fmt.Println("   go run main.go                 # Start your API server")
}
