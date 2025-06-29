package generator

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
)

// GenerateMainFile generates the main.go file that serves as the entrypoint
func GenerateMainFile(spec *models.OpenAPISpec, baseDir string, moduleName string) error {
	mainTemplate := `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ModuleName}}/api"
	"{{.ModuleName}}/server"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create API implementation
	apiImpl := api.NewAPI()

	// Create server with options
	srv := server.NewServer(apiImpl,
		server.WithMode("debug"), // Change to "release" for production
	)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: srv.GetRouter(),
	}

	// Setup graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancelShutdown := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancelShutdown()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		serverStopCtx()
	}()

	// Start server
	log.Printf("Starting server on :%s", port)
	log.Printf("API endpoints:")
	{{range .Endpoints}}
	log.Printf("  {{.Method}} {{.Path}} - {{.Description}}")
	{{end}}
	log.Printf("Press Ctrl+C to stop the server")

	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Printf("Server stopped gracefully")
}
`

	// Prepare endpoint data for logging
	var endpoints []struct {
		Method      string
		Path        string
		Description string
	}

	for path, operations := range spec.Paths {
		for method, op := range operations {
			description := op.Summary
			if description == "" {
				description = op.Description
			}
			if description == "" {
				description = "No description available"
			}

			endpoints = append(endpoints, struct {
				Method      string
				Path        string
				Description string
			}{
				Method:      method,
				Path:        path,
				Description: description,
			})
		}
	}

	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}

	data := struct {
		ModuleName string
		Endpoints  []struct {
			Method      string
			Path        string
			Description string
		}
	}{
		ModuleName: moduleName,
		Endpoints:  endpoints,
	}

	f, err := os.Create(filepath.Join(baseDir, "main.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
