// Package events defines all event types for the DoT application.
package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/darian/dot/messages/types"
)

func init() {
	dogma.RegisterEvent[*DeviceRegistered]("d4e5f6a7-b8c9-0123-defa-234567890123")
	dogma.RegisterEvent[*DeviceStateChanged]("e5f6a7b8-c9d0-1234-efab-345678901234")
	dogma.RegisterEvent[*DeviceRetired]("f6a7b8c9-d0e1-2345-fabc-456789012345")
}

// DeviceRegistered is an event indicating that a new device has been registered.
type DeviceRegistered struct {
	DeviceID     string            `json:"device_id"`
	DeviceName   string            `json:"device_name"`
	DeviceType   types.DeviceType  `json:"device_type"`
	Protocol     types.Protocol    `json:"protocol"`
	RoomID       string              `json:"room_id,omitempty"`
	Manufacturer string              `json:"manufacturer,omitempty"`
	Model        string              `json:"model,omitempty"`
	ExternalID   string              `json:"external_id,omitempty"`
	Capabilities []string            `json:"capabilities,omitempty"`
	RegisteredAt time.Time           `json:"registered_at"`
}

// MessageDescription returns a human-readable description of the message.
func (m *DeviceRegistered) MessageDescription() string {
	return fmt.Sprintf("device %q (%s) registered", m.DeviceName, m.DeviceType)
}

// Validate returns a non-nil error if the message is invalid.
func (m *DeviceRegistered) Validate(_ dogma.EventValidationScope) error {
	if m.DeviceID == "" {
		return errors.New("DeviceRegistered must have a device ID")
	}
	if m.DeviceName == "" {
		return errors.New("DeviceRegistered must have a device name")
	}
	if m.DeviceType == "" {
		return errors.New("DeviceRegistered must have a device type")
	}
	if m.Protocol == "" {
		return errors.New("DeviceRegistered must have a protocol")
	}
	return nil
}

// MarshalBinary returns a binary representation of the message.
func (m *DeviceRegistered) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
func (m *DeviceRegistered) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// DeviceStateChanged is an event indicating that a device state has changed.
type DeviceStateChanged struct {
	DeviceID  string            `json:"device_id"`
	State     map[string]string `json:"state"`
	ChangedAt time.Time         `json:"changed_at"`
}

// MessageDescription returns a human-readable description of the message.
func (m *DeviceStateChanged) MessageDescription() string {
	return fmt.Sprintf("device %q state changed", m.DeviceID)
}

// Validate returns a non-nil error if the message is invalid.
func (m *DeviceStateChanged) Validate(_ dogma.EventValidationScope) error {
	if m.DeviceID == "" {
		return errors.New("DeviceStateChanged must have a device ID")
	}
	if len(m.State) == 0 {
		return errors.New("DeviceStateChanged must have at least one state attribute")
	}
	return nil
}

// MarshalBinary returns a binary representation of the message.
func (m *DeviceStateChanged) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
func (m *DeviceStateChanged) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// DeviceRetired is an event indicating that a device has been retired.
type DeviceRetired struct {
	DeviceID  string    `json:"device_id"`
	RetiredAt time.Time `json:"retired_at"`
}

// MessageDescription returns a human-readable description of the message.
func (m *DeviceRetired) MessageDescription() string {
	return fmt.Sprintf("device %q retired", m.DeviceID)
}

// Validate returns a non-nil error if the message is invalid.
func (m *DeviceRetired) Validate(_ dogma.EventValidationScope) error {
	if m.DeviceID == "" {
		return errors.New("DeviceRetired must have a device ID")
	}
	return nil
}

// MarshalBinary returns a binary representation of the message.
func (m *DeviceRetired) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
func (m *DeviceRetired) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
