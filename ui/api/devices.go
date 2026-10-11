// Package api implements HTTP handlers for the DoT API gateway.
package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/darian/dot/messages/commands"
	"github.com/darian/dot/messages/events"
	"github.com/darian/dot/messages/types"
	"github.com/dogmatiq/dogma"
	"github.com/google/uuid"
)

// Device represents a device as returned by the API.
type Device struct {
	DeviceID     string            `json:"device_id"`
	DeviceName   string            `json:"device_name"`
	DeviceType   string            `json:"device_type"`
	Protocol     string            `json:"protocol"`
	RoomID       *string           `json:"room_id,omitempty"`
	Manufacturer *string           `json:"manufacturer,omitempty"`
	Model        *string           `json:"model,omitempty"`
	ExternalID   *string           `json:"external_id,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	State        map[string]string `json:"state"`
	Status       string            `json:"status"`
	Retired      bool              `json:"retired"`
	RegisteredAt time.Time         `json:"registered_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// RegisterDeviceRequest is the JSON body for registering a device.
type RegisterDeviceRequest struct {
	DeviceName   string   `json:"device_name"`
	DeviceType   string   `json:"device_type"`
	Protocol     string   `json:"protocol"`
	RoomID       string   `json:"room_id,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	ExternalID   string   `json:"external_id,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

// Handler holds the dependencies for API handlers.
type Handler struct {
	CommandExecutor dogma.CommandExecutor
	ReadDB          *sql.DB
}

// RegisterDevice handles POST /api/devices — registers a new device.
func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deviceID := uuid.New().String()

	err := h.CommandExecutor.ExecuteCommand(
		r.Context(),
		&commands.RegisterDevice{
			DeviceID:     deviceID,
			DeviceName:   req.DeviceName,
			DeviceType:   types.DeviceType(req.DeviceType),
			Protocol:     types.Protocol(req.Protocol),
			RoomID:       req.RoomID,
			Manufacturer: req.Manufacturer,
			Model:        req.Model,
			ExternalID:   req.ExternalID,
			Capabilities: req.Capabilities,
		},
		dogma.WithEventObserver(func(ctx context.Context, e *events.DeviceRegistered) (bool, error) {
			return true, nil
		}),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to register device: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"device_id": deviceID,
		"message":   "device registered",
	})
}

// ListDevices handles GET /api/devices — lists all devices from the projection.
func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	rows, err := h.ReadDB.QueryContext(
		r.Context(),
		`SELECT device_id, device_name, device_type, protocol, room_id,
			manufacturer, model, external_id, capabilities, state,
			status, retired, registered_at, updated_at
		FROM device_states
		WHERE retired = 0
		ORDER BY device_name`,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		var capsSQL, stateSQL string
		var retired int

		err := rows.Scan(
			&d.DeviceID, &d.DeviceName, &d.DeviceType, &d.Protocol,
			&d.RoomID, &d.Manufacturer, &d.Model, &d.ExternalID,
			&capsSQL, &stateSQL, &d.Status, &retired,
			&d.RegisteredAt, &d.UpdatedAt,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("scan failed: %v", err), http.StatusInternalServerError)
			return
		}

		d.Retired = retired == 1
		json.Unmarshal([]byte(capsSQL), &d.Capabilities)
		json.Unmarshal([]byte(stateSQL), &d.State)
		devices = append(devices, d)
	}

	if devices == nil {
		devices = []Device{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// GetDevice handles GET /api/devices/:id — returns a single device.
func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := r.PathValue("id")

	var d Device
	var capsSQL, stateSQL string
	var retired int

	err := h.ReadDB.QueryRowContext(
		r.Context(),
		`SELECT device_id, device_name, device_type, protocol, room_id,
			manufacturer, model, external_id, capabilities, state,
			status, retired, registered_at, updated_at
		FROM device_states WHERE device_id = ?`,
		deviceID,
	).Scan(
		&d.DeviceID, &d.DeviceName, &d.DeviceType, &d.Protocol,
		&d.RoomID, &d.Manufacturer, &d.Model, &d.ExternalID,
		&capsSQL, &stateSQL, &d.Status, &retired,
		&d.RegisteredAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("query failed: %v", err), http.StatusInternalServerError)
		return
	}

	d.Retired = retired == 1
	json.Unmarshal([]byte(capsSQL), &d.Capabilities)
	json.Unmarshal([]byte(stateSQL), &d.State)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// SetupRoutes registers all API routes.
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/devices", h.RegisterDevice)
	mux.HandleFunc("GET /api/devices", h.ListDevices)
	mux.HandleFunc("GET /api/devices/{id}", h.GetDevice)
}
