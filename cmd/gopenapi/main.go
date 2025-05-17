// main.go - Command-line interface for gopenapi
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shubhamku044/gopenapi/internal/generator"
	"github.com/shubhamku044/gopenapi/internal/output"
	"github.com/shubhamku044/gopenapi/internal/parser"
)

const (
	version = "0.1.0"
)

// convertOutputFiles converts generator.OutputFile slice to output.OutputFile slice
func convertOutputFiles(files []generator.OutputFile) []output.OutputFile {
	outputFiles := make([]output.OutputFile, len(files))
	for i, file := range files {
		outputFiles[i] = output.OutputFile{
			Path:    file.Path,
			Content: file.Content,
		}
	}
	return outputFiles
}

func main() {
	// Define command-line flags
	var (
		inputFile    string
		outputDir    string
		packageName  string
		serverOnly   bool
		clientOnly   bool
		templatesDir string
		showVersion  bool
	)

	flag.StringVar(&inputFile, "input", "", "Input OpenAPI YAML file (required)")
	flag.StringVar(&outputDir, "output", "./out", "Output directory for generated code")
	flag.StringVar(&packageName, "package", "", "Package name for generated code (default: derived from output dir)")
	flag.BoolVar(&serverOnly, "server-only", false, "Generate server code only")
	flag.BoolVar(&clientOnly, "client-only", false, "Generate client code only")
	flag.StringVar(&templatesDir, "templates", "", "Directory containing custom templates")
	flag.BoolVar(&showVersion, "version", false, "Show version information")

	// Parse the flags
	flag.Parse()

	// Show version and exit if requested
	if showVersion {
		fmt.Printf("gopenapi version %s\n", version)
		os.Exit(0)
	}

	// Check for required flags
	if inputFile == "" {
		fmt.Println("Error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

	// If package name is not specified, derive it from output directory
	if packageName == "" {
		packageName = filepath.Base(outputDir)
	}

	// Ensure input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: input file %s does not exist\n", inputFile)
		os.Exit(1)
	}

	// Parse the OpenAPI specification
	fmt.Println("Parsing OpenAPI specification...")
	spec, err := parser.ParseFile(inputFile)
	if err != nil {
		fmt.Printf("Error parsing specification: %v\n", err)
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Configure the generator
	genConfig := generator.Config{
		PackageName:  packageName,
		ServerOnly:   serverOnly,
		ClientOnly:   clientOnly,
		TemplatesDir: templatesDir,
	}

	// Generate the code
	fmt.Println("Generating code...")
	files, err := generator.Generate(spec, genConfig)
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}

	// Write files to disk
	fmt.Println("Writing files...")
	writer := output.NewWriter(outputDir)
	
	// Convert generator.OutputFile to output.OutputFile
	outputFiles := convertOutputFiles(files)
	
	if err := writer.WriteFiles(outputFiles); err != nil {
		fmt.Printf("Error writing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Code generation complete!")
}
