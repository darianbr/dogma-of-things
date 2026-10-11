package projections

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/darian/dot/messages/events"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// DeviceStateProjectionHandler maintains the current state of all devices.
type DeviceStateProjectionHandler struct {
	sqlprojection.NoCompactBehavior
}

// Configure configures the engine for this projection.
func (h *DeviceStateProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("device-states", "22222222-2222-2222-2222-222222222222")

	c.Routes(
		dogma.HandlesEvent[*events.DeviceRegistered](),
		dogma.HandlesEvent[*events.DeviceStateChanged](),
		dogma.HandlesEvent[*events.DeviceRetired](),
	)
}

// HandleEvent updates the device_states table when device events occur.
func (h *DeviceStateProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	m dogma.Event,
) error {
	switch x := m.(type) {
	case *events.DeviceRegistered:
		return h.deviceRegistered(ctx, tx, x)
	case *events.DeviceStateChanged:
		return h.deviceStateChanged(ctx, tx, x)
	case *events.DeviceRetired:
		return h.deviceRetired(ctx, tx, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func (h *DeviceStateProjectionHandler) deviceRegistered(
	ctx context.Context,
	tx *sql.Tx,
	x *events.DeviceRegistered,
) error {
	caps, _ := json.Marshal(x.Capabilities)

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO device_states (
			device_id, device_name, device_type, protocol, room_id,
			manufacturer, model, external_id, capabilities, state,
			status, retired, registered_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		x.DeviceID, x.DeviceName, string(x.DeviceType), string(x.Protocol),
		sql.NullString{String: x.RoomID, Valid: x.RoomID != ""},
		sql.NullString{String: x.Manufacturer, Valid: x.Manufacturer != ""},
		sql.NullString{String: x.Model, Valid: x.Model != ""},
		sql.NullString{String: x.ExternalID, Valid: x.ExternalID != ""},
		string(caps),
		"{}",
		"unknown",
		x.RegisteredAt.Format(time.RFC3339),
		time.Now().UTC().Format(time.RFC3339),
	)
	return err
}

func (h *DeviceStateProjectionHandler) deviceStateChanged(
	ctx context.Context,
	tx *sql.Tx,
	x *events.DeviceStateChanged,
) error {
	stateJSON, _ := json.Marshal(x.State)

	_, err := tx.ExecContext(
		ctx,
		`UPDATE device_states SET state = ?, updated_at = ? WHERE device_id = ?`,
		string(stateJSON),
		time.Now().UTC().Format(time.RFC3339),
		x.DeviceID,
	)
	return err
}

func (h *DeviceStateProjectionHandler) deviceRetired(
	ctx context.Context,
	tx *sql.Tx,
	x *events.DeviceRetired,
) error {
	_, err := tx.ExecContext(
		ctx,
		`UPDATE device_states SET retired = 1, updated_at = ? WHERE device_id = ?`,
		time.Now().UTC().Format(time.RFC3339),
		x.DeviceID,
	)
	return err
}

// Reset clears all projection data.
func (h *DeviceStateProjectionHandler) Reset(
	ctx context.Context,
	tx *sql.Tx,
	_ dogma.ProjectionResetScope,
) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM device_states`)
	return err
}
