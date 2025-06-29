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

	spec, err := parser.ParseSpecFile(*specFile)
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI specification: %v", err)
	}

	pkg := *packageName
	if pkg == "" {
		absPath, err := filepath.Abs(*outputDir)
		if err == nil {
			pkg = filepath.Base(absPath)
		}

		if pkg == "" || pkg == "." {
			pkg = "generated"
		}
	}

	moduleName := pkg

	config := generator.Config{
		OutputDir:   *outputDir,
		PackageName: pkg,
		ModuleName:  moduleName,
	}

	err = generator.GenerateCode(spec, config)
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	fmt.Printf("Code successfully generated in %s\n", *outputDir)
}
