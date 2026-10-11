// Package commands defines all command types for the DoT application.
package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/darian/dot/messages/types"
)

func init() {
	dogma.RegisterCommand[*RegisterDevice]("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	dogma.RegisterCommand[*SetDeviceState]("b2c3d4e5-f6a7-8901-bcde-f12345678901")
	dogma.RegisterCommand[*RetireDevice]("c3d4e5f6-a7b8-9012-cdef-123456789012")
}

// RegisterDevice commands the system to register a new IoT device.
type RegisterDevice struct {
	DeviceID      string            `json:"device_id"`
	DeviceName    string            `json:"device_name"`
	DeviceType    types.DeviceType  `json:"device_type"`
	Protocol      types.Protocol    `json:"protocol"`
	RoomID        string             `json:"room_id,omitempty"`
	Manufacturer  string             `json:"manufacturer,omitempty"`
	Model         string             `json:"model,omitempty"`
	ExternalID    string             `json:"external_id,omitempty"`
	Capabilities  []string           `json:"capabilities,omitempty"`
}

// MessageDescription returns a human-readable description of the message.
func (m *RegisterDevice) MessageDescription() string {
	return fmt.Sprintf("register device %q (%s)", m.DeviceName, m.DeviceType)
}

// Validate returns a non-nil error if the message is invalid.
func (m *RegisterDevice) Validate(_ dogma.CommandValidationScope) error {
	if m.DeviceID == "" {
		return errors.New("RegisterDevice must have a device ID")
	}
	if m.DeviceName == "" {
		return errors.New("RegisterDevice must have a device name")
	}
	if m.DeviceType == "" {
		return errors.New("RegisterDevice must have a device type")
	}
	if m.Protocol == "" {
		return errors.New("RegisterDevice must have a protocol")
	}
	return nil
}

// MarshalBinary returns a binary representation of the message.
func (m *RegisterDevice) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
func (m *RegisterDevice) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// SetDeviceState commands a device to change its state.
type SetDeviceState struct {
	DeviceID  string            `json:"device_id"`
	State     map[string]string `json:"state"`
	Timestamp time.Time         `json:"timestamp,omitempty"`
}

// MessageDescription returns a human-readable description of the message.
func (m *SetDeviceState) MessageDescription() string {
	return fmt.Sprintf("set state of device %q", m.DeviceID)
}

// Validate returns a non-nil error if the message is invalid.
func (m *SetDeviceState) Validate(_ dogma.CommandValidationScope) error {
	if m.DeviceID == "" {
		return errors.New("SetDeviceState must have a device ID")
	}
	if len(m.State) == 0 {
		return errors.New("SetDeviceState must have at least one state attribute")
	}
	return nil
}

// MarshalBinary returns a binary representation of the message.
func (m *SetDeviceState) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
func (m *SetDeviceState) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// RetireDevice commands the system to retire a device.
type RetireDevice struct {
	DeviceID string `json:"device_id"`
}

// MessageDescription returns a human-readable description of the message.
func (m *RetireDevice) MessageDescription() string {
	return fmt.Sprintf("retire device %q", m.DeviceID)
}

// Validate returns a non-nil error if the message is invalid.
func (m *RetireDevice) Validate(_ dogma.CommandValidationScope) error {
	if m.DeviceID == "" {
		return errors.New("RetireDevice must have a device ID")
	}
	return nil
}

// MarshalBinary returns a binary representation of the message.
func (m *RetireDevice) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
func (m *RetireDevice) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
