// Package webui provides an embedded web server for managing aliases
// through a browser-based user interface.
package webui

import (
	"io/fs"
	"net/http"

	"aliasly/web"
)

// Server represents the web UI server.
// It handles HTTP requests and serves both the static files
// and the API endpoints.
type Server struct {
	// mux is the HTTP request multiplexer (router)
	// It routes incoming requests to the appropriate handlers
	mux *http.ServeMux
}

// NewServer creates a new web UI server instance.
// It sets up all routes and handlers.
func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}

	// Set up routes
	s.setupRoutes()

	return s
}

// Handler returns the HTTP handler for this server.
// This is used by the http.Server to handle incoming requests.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// setupRoutes configures all the URL routes for the server.
func (s *Server) setupRoutes() {
	// API routes for CRUD operations on aliases
	// These return JSON and are called by the JavaScript frontend

	// GET /api/aliases - List all aliases
	s.mux.HandleFunc("GET /api/aliases", handleListAliases)

	// POST /api/aliases - Create a new alias
	s.mux.HandleFunc("POST /api/aliases", handleCreateAlias)

	// PUT /api/aliases/{name} - Update an existing alias
	s.mux.HandleFunc("PUT /api/aliases/{name}", handleUpdateAlias)

	// DELETE /api/aliases/{name} - Delete an alias
	s.mux.HandleFunc("DELETE /api/aliases/{name}", handleDeleteAlias)

	// Serve static files (HTML, CSS, JS)
	// We need to strip the "static" prefix because the files are
	// embedded under "static/" but we want to serve them from "/"
	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		// This should never happen since we control the embed directive
		panic("failed to get static files: " + err.Error())
	}

	// http.FileServer creates a handler that serves files from the filesystem
	// We wrap it to serve index.html for the root path
	fileServer := http.FileServer(http.FS(staticFS))
	s.mux.Handle("/", fileServer)
}
