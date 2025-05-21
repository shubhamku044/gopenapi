package generator

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
)

// GenerateReadme generates a README file for the generated API
func GenerateReadme(spec *models.OpenAPISpec, baseDir string, packageName string) error {
	readmeTemplate := `# Generated API for {{.Title}}

This code was generated using GopenAPI for the {{.Title}} API (version {{.Version}}).

## Usage

To use this generated API in your project:

### 1. Implement the API interface

Create a custom implementation of the API interface:

` + "```go" + `
package main

import (
	"github.com/gin-gonic/gin"
	
	"{{.PackageName}}/generated/api"
	"{{.PackageName}}/generated/server"
)

// CustomAPI implements the generated API interface
type CustomAPI struct {
	// Add your dependencies here (database connections, etc.)
}

// Ensure CustomAPI implements the API interface
var _ api.API = (*CustomAPI)(nil)

// Example implementation of a generated method
func (a *CustomAPI) ListUsers(c *gin.Context) {
	// Your implementation here
	users := []map[string]interface{}{
		{"id": "1", "name": "John Doe"},
		{"id": "2", "name": "Jane Smith"},
	}
	c.JSON(200, users)
}

// Implement all other methods defined in the API interface...
` + "```" + `

### 2. Create and run the server

` + "```go" + `
package main

import (
	"log"
	
	"{{.PackageName}}/generated/api"
	"{{.PackageName}}/generated/server"
)

func main() {
	// Create your API implementation
	customAPI := &CustomAPI{}
	
	// Create the server with your API implementation
	srv := server.NewServer(customAPI)
	
	// Start the server
	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
` + "```" + `

### 3. Or use with an existing Gin application

` + "```go" + `
package main

import (
	"github.com/gin-gonic/gin"
	
	"{{.PackageName}}/generated/api"
	"{{.PackageName}}/generated/server"
)

func main() {
	// Create your existing Gin router
	router := gin.Default()
	
	// Add some custom middleware or routes
	router.Use(customMiddleware())
	router.GET("/health", healthCheckHandler)
	
	// Create your API implementation
	customAPI := &CustomAPI{}
	
	// Create the server with your API implementation
	srv := server.NewServer(customAPI, server.WithMiddleware(loggerMiddleware()))
	
	// Get the Gin router from the server and add it to your existing router
	apiRouter := srv.GetRouter()
	
	// Use the API routes in your application
	router.Any("/api/*path", func(c *gin.Context) {
		apiRouter.HandleContext(c)
	})
	
	// Start your server
	router.Run(":8080")
}
` + "```" + `
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return err
	}

	data := struct {
		Title       string
		Version     string
		PackageName string
	}{
		Title:       spec.Info.Title,
		Version:     spec.Info.Version,
		PackageName: packageName,
	}

	f, err := os.Create(filepath.Join(baseDir, "README.md"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
