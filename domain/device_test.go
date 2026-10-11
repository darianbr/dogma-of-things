package domain_test

import (
	"errors"
	"testing"

	dot "github.com/darian/dot"
	"github.com/darian/dot/messages/commands"
	"github.com/darian/dot/messages/events"
	"github.com/darian/dot/messages/types"
	. "github.com/dogmatiq/testkit"
)

// errMismatch is returned by predicate functions when a recorded event does not
// match the expected values. It is used with ToRecordEventMatching, which is
// required here because the device events carry non-deterministic timestamps
// (set via time.Now() in the handlers) that preclude exact ToRecordEvent matches.
var errMismatch = errors.New("event did not match expected values")

func Test_RegisterDevice(t *testing.T) {
	t.Run(
		"when the device does not exist",
		func(t *testing.T) {
			t.Run(
				"the device is registered",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Expect(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:     "dev-001",
									DeviceName:   "Living Room Light",
									DeviceType:   types.DeviceTypeLight,
									Protocol:     types.ProtocolWiFi,
									RoomID:       "room-living",
									Manufacturer: "Eufy",
									Model:        "E22",
									ExternalID:   "light.living_room_main",
									Capabilities: []string{"onoff", "brightness", "color_temp"},
								},
							),
							ToRecordEventMatching(func(e *events.DeviceRegistered) error {
								if e.DeviceID != "dev-001" ||
									e.DeviceName != "Living Room Light" ||
									e.DeviceType != types.DeviceTypeLight ||
									e.Protocol != types.ProtocolWiFi ||
									e.RoomID != "room-living" ||
									e.Manufacturer != "Eufy" ||
									e.Model != "E22" ||
									e.ExternalID != "light.living_room_main" {
									return errMismatch
								}
								if len(e.Capabilities) != 3 ||
									e.Capabilities[0] != "onoff" ||
									e.Capabilities[1] != "brightness" ||
									e.Capabilities[2] != "color_temp" {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)

			t.Run(
				"the device is registered with minimal fields",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Expect(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-002",
									DeviceName: "Bedroom Sensor",
									DeviceType: types.DeviceTypeSensor,
									Protocol:   types.ProtocolZigbee,
								},
							),
							ToRecordEventMatching(func(e *events.DeviceRegistered) error {
								if e.DeviceID != "dev-002" ||
									e.DeviceName != "Bedroom Sensor" ||
									e.DeviceType != types.DeviceTypeSensor ||
									e.Protocol != types.ProtocolZigbee ||
									e.RoomID != "" ||
									e.Manufacturer != "" ||
									e.Model != "" ||
									e.ExternalID != "" {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)
		},
	)

	t.Run(
		"when the device already exists",
		func(t *testing.T) {
			t.Run(
				"nothing happens to the existing device",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-001",
									DeviceName: "Living Room Light",
									DeviceType: types.DeviceTypeLight,
									Protocol:   types.ProtocolWiFi,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-001",
									DeviceName: "Living Room Light",
									DeviceType: types.DeviceTypeLight,
									Protocol:   types.ProtocolWiFi,
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.DeviceRegistered{}),
							),
						)
				},
			)
		},
	)
}

func Test_SetDeviceState(t *testing.T) {
	t.Run(
		"when the device is registered",
		func(t *testing.T) {
			t.Run(
				"the device state is updated",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-003",
									DeviceName: "Kitchen Switch",
									DeviceType: types.DeviceTypeSwitch,
									Protocol:   types.ProtocolZWave,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-003",
									State: map[string]string{
										"on":         "true",
										"brightness": "80",
									},
								},
							),
							ToRecordEventMatching(func(e *events.DeviceStateChanged) error {
								if e.DeviceID != "dev-003" ||
									e.State["on"] != "true" ||
									e.State["brightness"] != "80" {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)

			t.Run(
				"the previous state is overwritten",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-004",
									DeviceName: "Hallway Light",
									DeviceType: types.DeviceTypeLight,
									Protocol:   types.ProtocolWiFi,
								},
							),
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-004",
									State: map[string]string{
										"on":         "true",
										"brightness": "100",
									},
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-004",
									State: map[string]string{
										"on":         "false",
										"brightness": "0",
									},
								},
							),
							ToRecordEventMatching(func(e *events.DeviceStateChanged) error {
								if e.State["on"] != "false" || e.State["brightness"] != "0" {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)
		},
	)

	t.Run(
		"when the device is not registered",
		func(t *testing.T) {
			t.Run(
				"the state change is rejected",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-nonexistent",
									State: map[string]string{
										"on": "true",
									},
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.DeviceStateChanged{}),
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the device is retired",
		func(t *testing.T) {
			t.Run(
				"the state change is rejected",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-007",
									DeviceName: "Retired Light",
									DeviceType: types.DeviceTypeLight,
									Protocol:   types.ProtocolWiFi,
								},
							),
							ExecuteCommand(
								&commands.RetireDevice{
									DeviceID: "dev-007",
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-007",
									State: map[string]string{
										"on": "true",
									},
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.DeviceStateChanged{}),
							),
						)
				},
			)
		},
	)
}

func Test_RetireDevice(t *testing.T) {
	t.Run(
		"when the device is registered",
		func(t *testing.T) {
			t.Run(
				"the device is retired",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-005",
									DeviceName: "Old Thermostat",
									DeviceType: types.DeviceTypeThermostat,
									Protocol:   types.ProtocolZigbee,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.RetireDevice{
									DeviceID: "dev-005",
								},
							),
							ToRecordEventMatching(func(e *events.DeviceRetired) error {
								if e.DeviceID != "dev-005" {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)
		},
	)

	t.Run(
		"when the device is not registered",
		func(t *testing.T) {
			t.Run(
				"the retirement is rejected",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Expect(
							ExecuteCommand(
								&commands.RetireDevice{
									DeviceID: "dev-nonexistent",
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.DeviceRetired{}),
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the device is already retired",
		func(t *testing.T) {
			t.Run(
				"nothing happens to the existing device",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-006",
									DeviceName: "Garage Door",
									DeviceType: types.DeviceTypeSwitch,
									Protocol:   types.ProtocolWiFi,
								},
							),
							ExecuteCommand(
								&commands.RetireDevice{
									DeviceID: "dev-006",
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.RetireDevice{
									DeviceID: "dev-006",
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.DeviceRetired{}),
							),
						)
				},
			)
		},
	)
}

func Test_DeviceLifecycle(t *testing.T) {
	t.Run(
		"when a device is registered, used, and retired",
		func(t *testing.T) {
			t.Run(
				"the full lifecycle is recorded",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Expect(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:     "dev-lifecycle",
									DeviceName:   "Master Bedroom Light",
									DeviceType:   types.DeviceTypeLight,
									Protocol:     types.ProtocolZigbee,
									RoomID:       "room-master",
									Manufacturer: "Eufy",
									Model:        "E22",
									Capabilities: []string{"onoff", "brightness"},
								},
							),
							ToRecordEventOfType(&events.DeviceRegistered{}),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-lifecycle",
									State: map[string]string{
										"on":         "true",
										"brightness": "100",
									},
								},
							),
							ToRecordEventOfType(&events.DeviceStateChanged{}),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-lifecycle",
									State: map[string]string{
										"on":         "true",
										"brightness": "50",
									},
								},
							),
							ToRecordEventOfType(&events.DeviceStateChanged{}),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-lifecycle",
									State: map[string]string{
										"on":         "false",
										"brightness": "0",
									},
								},
							),
							ToRecordEventOfType(&events.DeviceStateChanged{}),
						).
						Expect(
							ExecuteCommand(
								&commands.RetireDevice{
									DeviceID: "dev-lifecycle",
								},
							),
							ToRecordEventOfType(&events.DeviceRetired{}),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-lifecycle",
									State: map[string]string{
										"on": "true",
									},
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.DeviceStateChanged{}),
							),
						)
				},
			)
		},
	)
}

func Test_EventTimestamps(t *testing.T) {
	t.Run(
		"when a device is registered",
		func(t *testing.T) {
			t.Run(
				"the event has a timestamp",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Expect(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-time",
									DeviceName: "Timestamp Test",
									DeviceType: types.DeviceTypeSensor,
									Protocol:   types.ProtocolWiFi,
								},
							),
							ToRecordEventMatching(func(e *events.DeviceRegistered) error {
								if e.RegisteredAt.IsZero() {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)
		},
	)

	t.Run(
		"when a device state changes",
		func(t *testing.T) {
			t.Run(
				"the event has a timestamp",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-time-2",
									DeviceName: "State Timestamp Test",
									DeviceType: types.DeviceTypeLight,
									Protocol:   types.ProtocolWiFi,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.SetDeviceState{
									DeviceID: "dev-time-2",
									State: map[string]string{
										"on": "true",
									},
								},
							),
							ToRecordEventMatching(func(e *events.DeviceStateChanged) error {
								if e.ChangedAt.IsZero() {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)
		},
	)

	t.Run(
		"when a device is retired",
		func(t *testing.T) {
			t.Run(
				"the event has a timestamp",
				func(t *testing.T) {
					Begin(t, &dot.App{}).
						Prepare(
							ExecuteCommand(
								&commands.RegisterDevice{
									DeviceID:   "dev-time-3",
									DeviceName: "Retire Timestamp Test",
									DeviceType: types.DeviceTypeSwitch,
									Protocol:   types.ProtocolWiFi,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.RetireDevice{
									DeviceID: "dev-time-3",
								},
							),
							ToRecordEventMatching(func(e *events.DeviceRetired) error {
								if e.RetiredAt.IsZero() {
									return errMismatch
								}
								return nil
							}),
						)
				},
			)
		},
	)
}
