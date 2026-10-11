// Package domain implements the business logic for the DoT application.
package domain

import (
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/darian/dot/messages/commands"
	"github.com/darian/dot/messages/events"
	"github.com/darian/dot/messages/types"
)

// device is the aggregate root for an IoT device.
type device struct {
	dogma.NoSnapshotBehavior

	DeviceID     string
	DeviceName   string
	DeviceType   types.DeviceType
	Protocol     types.Protocol
	RoomID       string
	Manufacturer string
	Model        string
	ExternalID   string
	Capabilities []string
	State        map[string]string
	Status       types.DeviceStatus
	Retired      bool
}

func (d *device) AggregateInstanceDescription() string {
	if d.DeviceName == "" {
		return ""
	}
	if d.Retired {
		return fmt.Sprintf("%s (retired)", d.DeviceName)
	}
	return fmt.Sprintf("%s (%s, %s)", d.DeviceName, d.DeviceType, d.Status)
}

// ApplyEvent rebuilds aggregate state from an event.
func (d *device) ApplyEvent(m dogma.Event) {
	switch m := m.(type) {
	case *events.DeviceRegistered:
		d.DeviceID = m.DeviceID
		d.DeviceName = m.DeviceName
		d.DeviceType = m.DeviceType
		d.Protocol = m.Protocol
		d.RoomID = m.RoomID
		d.Manufacturer = m.Manufacturer
		d.Model = m.Model
		d.ExternalID = m.ExternalID
		d.Capabilities = m.Capabilities
		d.Status = types.DeviceStatusUnknown

	case *events.DeviceStateChanged:
		d.State = m.State

	case *events.DeviceRetired:
		d.Retired = true
	}
}

// RegisterDevice handles the RegisterDevice command.
func (d *device) RegisterDevice(s dogma.AggregateCommandScope[*device], m *commands.RegisterDevice) {
	if d.DeviceName != "" {
		s.Log("device already registered")
		return
	}

	s.RecordEvent(&events.DeviceRegistered{
		DeviceID:     m.DeviceID,
		DeviceName:   m.DeviceName,
		DeviceType:   m.DeviceType,
		Protocol:     m.Protocol,
		RoomID:       m.RoomID,
		Manufacturer: m.Manufacturer,
		Model:        m.Model,
		ExternalID:   m.ExternalID,
		Capabilities: m.Capabilities,
		RegisteredAt: time.Now(),
	})
}

// SetDeviceState handles the SetDeviceState command.
func (d *device) SetDeviceState(s dogma.AggregateCommandScope[*device], m *commands.SetDeviceState) {
	if d.DeviceName == "" {
		s.Log("device not registered")
		return
	}
	if d.Retired {
		s.Log("device is retired")
		return
	}

	s.RecordEvent(&events.DeviceStateChanged{
		DeviceID:  m.DeviceID,
		State:     m.State,
		ChangedAt: time.Now(),
	})
}

// RetireDevice handles the RetireDevice command.
func (d *device) RetireDevice(s dogma.AggregateCommandScope[*device], m *commands.RetireDevice) {
	if d.DeviceName == "" {
		s.Log("device not registered")
		return
	}
	if d.Retired {
		s.Log("device already retired")
		return
	}

	s.RecordEvent(&events.DeviceRetired{
		DeviceID:  m.DeviceID,
		RetiredAt: time.Now(),
	})
}

// DeviceHandler implements the business logic for IoT devices.
type DeviceHandler struct{}

// Configure configures the behavior of the engine as it relates to this handler.
func (DeviceHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("device", "11111111-1111-1111-1111-111111111111")

	c.Routes(
		dogma.HandlesCommand[*commands.RegisterDevice](),
		dogma.HandlesCommand[*commands.SetDeviceState](),
		dogma.HandlesCommand[*commands.RetireDevice](),
		dogma.RecordsEvent[*events.DeviceRegistered](),
		dogma.RecordsEvent[*events.DeviceStateChanged](),
		dogma.RecordsEvent[*events.DeviceRetired](),
	)
}

// New returns a new device instance.
func (DeviceHandler) New() *device {
	return &device{
		State: make(map[string]string),
	}
}

// RouteCommandToInstance returns the ID of the aggregate instance that is targeted by m.
func (DeviceHandler) RouteCommandToInstance(m dogma.Command) string {
	switch x := m.(type) {
	case *commands.RegisterDevice:
		return x.DeviceID
	case *commands.SetDeviceState:
		return x.DeviceID
	case *commands.RetireDevice:
		return x.DeviceID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (DeviceHandler) HandleCommand(
	d *device,
	s dogma.AggregateCommandScope[*device],
	m dogma.Command,
) {
	switch m := m.(type) {
	case *commands.RegisterDevice:
		d.RegisterDevice(s, m)
	case *commands.SetDeviceState:
		d.SetDeviceState(s, m)
	case *commands.RetireDevice:
		d.RetireDevice(s, m)
	default:
		panic(dogma.UnexpectedMessage)
	}
}
