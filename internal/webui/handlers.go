package webui

import (
	"encoding/json"
	"net/http"

	"aliasly/internal/alias"
	"aliasly/internal/config"
)

// APIResponse is a standard response format for our API.
// All API responses follow this structure for consistency.
type APIResponse struct {
	// Success indicates whether the operation succeeded
	Success bool `json:"success"`

	// Data contains the response data (if successful)
	Data interface{} `json:"data,omitempty"`

	// Error contains the error message (if failed)
	Error string `json:"error,omitempty"`
}

// handleListAliases handles GET /api/aliases
// It returns a list of all configured aliases as JSON.
func handleListAliases(w http.ResponseWriter, r *http.Request) {
	// Get all aliases from config
	aliases, err := alias.GetAll()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send success response with aliases
	sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    aliases,
	})
}

// handleCreateAlias handles POST /api/aliases
// It creates a new alias from the JSON request body.
func handleCreateAlias(w http.ResponseWriter, r *http.Request) {
	// Parse the request body as JSON
	var newAlias config.Alias
	if err := json.NewDecoder(r.Body).Decode(&newAlias); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if newAlias.Name == "" {
		sendError(w, http.StatusBadRequest, "Alias name is required")
		return
	}
	if newAlias.Command == "" {
		sendError(w, http.StatusBadRequest, "Command is required")
		return
	}

	// Check if alias already exists
	if _, exists := alias.Find(newAlias.Name); exists {
		sendError(w, http.StatusConflict, "Alias '"+newAlias.Name+"' already exists")
		return
	}

	// Add the alias
	if err := alias.Add(newAlias); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the created alias
	sendJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    newAlias,
	})
}

// handleUpdateAlias handles PUT /api/aliases/{name}
// It updates an existing alias with the JSON request body.
func handleUpdateAlias(w http.ResponseWriter, r *http.Request) {
	// Get the alias name from the URL path
	// In Go 1.22+, we can use PathValue to get path parameters
	aliasName := r.PathValue("name")
	if aliasName == "" {
		sendError(w, http.StatusBadRequest, "Alias name is required in URL")
		return
	}

	// Check if alias exists
	if _, exists := alias.Find(aliasName); !exists {
		sendError(w, http.StatusNotFound, "Alias '"+aliasName+"' not found")
		return
	}

	// Parse the request body
	var updatedAlias config.Alias
	if err := json.NewDecoder(r.Body).Decode(&updatedAlias); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Ensure the name matches the URL
	updatedAlias.Name = aliasName

	// Validate required fields
	if updatedAlias.Command == "" {
		sendError(w, http.StatusBadRequest, "Command is required")
		return
	}

	// Update the alias
	if err := alias.Update(updatedAlias); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the updated alias
	sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    updatedAlias,
	})
}

// handleDeleteAlias handles DELETE /api/aliases/{name}
// It deletes an existing alias.
func handleDeleteAlias(w http.ResponseWriter, r *http.Request) {
	// Get the alias name from the URL path
	aliasName := r.PathValue("name")
	if aliasName == "" {
		sendError(w, http.StatusBadRequest, "Alias name is required in URL")
		return
	}

	// Check if alias exists
	if _, exists := alias.Find(aliasName); !exists {
		sendError(w, http.StatusNotFound, "Alias '"+aliasName+"' not found")
		return
	}

	// Delete the alias
	if err := alias.Remove(aliasName); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return success
	sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
	})
}

// sendJSON sends a JSON response with the given status code.
// This is a helper function to avoid repeating JSON encoding code.
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	// Set the content type header before writing the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode the data as JSON and write to the response
	// json.NewEncoder writes directly to the http.ResponseWriter
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, log it (in production, use proper logging)
		// We can't change the status code at this point
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// sendError sends an error response with the given status code and message.
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}
