// Package mqtt is the main encompasing mqtt package
package mqtt

// ControlType packets for mqtt
const (
	reserved = iota
	Connect
	ConnectAck
	Publish
	PublishAck
	PublishReceived
	PublishRelease
	PublishComplete
	Subscribe
	SubscribeAck
	UnSubscribe
	UnSubscribeAck
	PingRequest
	PingResponse
	Disconnect
)

const (
	connectMinimumLength = 10
)
