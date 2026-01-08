package handlers

import (
	"net/http"

	"github.com/ajmasia/shellify/internal/domain"
	"github.com/ajmasia/shellify/internal/infrastructure/storage"
)

// SettingsHandler handles settings-related HTTP requests.
type SettingsHandler struct {
	storage *storage.Storage
}

// NewSettingsHandler creates a new SettingsHandler.
func NewSettingsHandler(store *storage.Storage) *SettingsHandler {
	return &SettingsHandler{
		storage: store,
	}
}

// Get returns the application settings.
func (h *SettingsHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config, err := h.storage.LoadConfig()
		if err != nil {
			// Return default config if none exists
			config = domain.DefaultConfig()
		}

		Success(w, config)
	}
}

// Update updates the application settings.
func (h *SettingsHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var config domain.Config
		if err := DecodeJSON(r, &config); err != nil {
			BadRequest(w, "invalid JSON")
			return
		}

		if err := h.storage.SaveConfig(config); err != nil {
			InternalError(w, err.Error())
			return
		}

		Success(w, config)
	}
}
