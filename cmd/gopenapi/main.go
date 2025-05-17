package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shubhamku044/gopenapi/pkg/generator"
)

func main() {
	// Check if any command was provided
	if len(os.Args) < 2 {
		fmt.Println("Expected 'generate' subcommand")
		os.Exit(1)
	}

	// Handle subcommands
	switch os.Args[1] {
	case "generate":
		handleGenerate(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func handleGenerate(args []string) {
	if len(args) < 1 {
		fmt.Println("Expected 'model', 'server', or 'all' subcommand for generate")
		os.Exit(1)
	}

	// Define command line flags
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	
	// Define flags for the generate command
	var inputFile string
	var outputDir string
	generateCmd.StringVar(&inputFile, "input", "", "Path to gopenapi spec file (YAML/JSON)")
	generateCmd.StringVar(&outputDir, "output", "./gen", "Output directory for generated code")

	// Parse flags after the subcommand
	generateCmd.Parse(args[1:])

	// Validate required flags
	if inputFile == "" {
		fmt.Println("--input flag is required")
		os.Exit(1)
	}

	// Ensure input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Input file does not exist: %s\n", inputFile)
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %s\n", err)
		os.Exit(1)
	}

	// Create a generator instance
	gen, err := generator.NewGenerator(inputFile, outputDir)
	if err != nil {
		fmt.Printf("Failed to create generator: %s\n", err)
		os.Exit(1)
	}

	// Load the gopenapi spec
	if err := gen.LoadSpec(); err != nil {
		fmt.Printf("Failed to load spec: %s\n", err)
		os.Exit(1)
	}

	// Handle different generation types
	switch args[0] {
	case "model":
		fmt.Printf("Generating models from %s to %s\n", inputFile, filepath.Join(outputDir, "models"))
		if err := gen.GenerateModels(); err != nil {
			fmt.Printf("Failed to generate models: %s\n", err)
			os.Exit(1)
		}
	case "server":
		fmt.Printf("Generating server code from %s to %s\n", inputFile, outputDir)
		if err := gen.GenerateServer(); err != nil {
			fmt.Printf("Failed to generate server: %s\n", err)
			os.Exit(1)
		}
	case "all":
		fmt.Printf("Generating all code from %s to %s\n", inputFile, outputDir)
		if err := gen.GenerateAll(); err != nil {
			fmt.Printf("Failed to generate all components: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown generate subcommand: %s\n", args[0])
		os.Exit(1)
	}

	fmt.Println("Code generation completed successfully!")
}
