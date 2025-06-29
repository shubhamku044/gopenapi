package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/shubhamku044/gopenapi/internal/generator"
	"github.com/shubhamku044/gopenapi/internal/parser"
)

func main() {
	specFile := flag.String("spec", "", "Path to OpenAPI specification file (YAML or JSON)")
	outputDir := flag.String("output", "./generated", "Output directory for generated code")
	packageName := flag.String("package", "", "Package name for generated code (optional)")
	flag.Parse()

	if *specFile == "" {
		log.Fatal("Please provide an OpenAPI specification file with --spec")
	}

	// Parse the OpenAPI specification
	spec, err := parser.ParseSpecFile(*specFile)
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI specification: %v", err)
	}

	// Determine package name
	pkg := *packageName
	if pkg == "" {
		// Try to detect from the output directory
		absPath, err := filepath.Abs(*outputDir)
		if err == nil {
			pkg = filepath.Base(filepath.Dir(absPath))
		}

		// If still empty, use a default
		if pkg == "" || pkg == "." {
			pkg = "gopenapi"
		}
	}

	// Generate the code
	config := generator.Config{
		OutputDir:   *outputDir,
		PackageName: pkg,
	}

	err = generator.GenerateCode(spec, config)
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	fmt.Printf("Code successfully generated in %s\n", *outputDir)
}
