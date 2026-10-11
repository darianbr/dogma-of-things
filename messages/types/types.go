// Package types defines shared types used across message definitions.
package types

// DeviceType categorizes the primary function of a device.
type DeviceType string

const (
	DeviceTypeLight        DeviceType = "light"
	DeviceTypeSwitch       DeviceType = "switch"
	DeviceTypeSensor       DeviceType = "sensor"
	DeviceTypeThermostat   DeviceType = "thermostat"
	DeviceTypeSpeaker      DeviceType = "speaker"
	DeviceTypeCamera       DeviceType = "camera"
	DeviceTypeVacuum       DeviceType = "vacuum"
	DeviceTypeLock         DeviceType = "lock"
	DeviceTypeBlind        DeviceType = "blind"
	DeviceTypeCover        DeviceType = "cover"
	DeviceTypeHumidifier   DeviceType = "humidifier"
	DeviceTypeMediaPlayer  DeviceType = "media_player"
	DeviceTypeClimate      DeviceType = "climate"
	DeviceTypeFan          DeviceType = "fan"
	DeviceTypeSiren        DeviceType = "siren"
	DeviceTypeAlarmControl DeviceType = "alarm_control_panel"
)

// Protocol identifies how the device communicates.
type Protocol string

const (
	ProtocolZigbee   Protocol = "zigbee"
	ProtocolZWave    Protocol = "zwave"
	ProtocolMQTT     Protocol = "mqtt"
	ProtocolWiFi     Protocol = "wifi"
	ProtocolBluetooth Protocol = "bluetooth"
	ProtocolZell     Protocol = "zell"
	ProtocolHomeKit  Protocol = "homekit"
	ProtocolOther    Protocol = "other"
)

// DeviceStatus represents the current connectivity state of a device.
type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
	DeviceStatusUnknown DeviceStatus = "unknown"
)
