package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed swagger-ui/*
var swaggerUI embed.FS

// Server represents the API server
type Server struct {
	echo      *echo.Echo
	dockerCli *client.Client
}

// New creates a new server instance
func New(dockerCli *client.Client) *Server {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve Swagger UI
	swaggerFS, err := fs.Sub(swaggerUI, "swagger-ui")
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFS)))))

	// Serve OpenAPI spec
	e.GET("/api/spec/project.yaml", func(c echo.Context) error {
		return c.File("api/spec/project.yaml")
	})

	return &Server{
		echo:      e,
		dockerCli: dockerCli,
	}
}

// Start starts the server
func (s *Server) Start(addr string) error {
	return s.echo.Start(addr)
}
