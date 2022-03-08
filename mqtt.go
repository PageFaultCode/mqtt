// Package mqtt is the main encompasing mqtt package
package mqtt

// MQTT is the main struct for MQTT messaging
type MQTT struct {
	isConnected bool
}

// Connect will connect to the given server
func (m *MQTT) Connect(url string, port uint16, useTLS bool) error {
	m.isConnected = true
	return nil
}

// Disconnect will disconnect from the server if connected
func (m *MQTT) Disconnect() error {
	if m.isConnected {
		m.isConnected = false
	}
	return nil
}
