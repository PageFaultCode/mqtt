// Package mqtt is the main encompasing mqtt package
package mqtt

import (
	"fmt"
	"net"
)

// MQTT is the main struct for MQTT messaging
type MQTT struct {
	isConnected bool
	connection  *net.Conn
}

// Connect will connect to the given server
func (m *MQTT) Connect(url string, port uint16, useTLS bool) error {
	targetURL := url + ":" + fmt.Sprintf("%d", port)
	connection, err := net.Dial("tcp", targetURL)
	if err != nil {
		return err
	}

	m.connection = &connection
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
