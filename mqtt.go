// Package mqtt is the main encompasing mqtt package
package mqtt

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"
)

const (
	defaultConnectTimeout = time.Second * 30
)

// MQTTListener is an interface defining
// what a listener of MQTT communications
// needs to implement
type MQTTListener interface {
	Identifier() string
	OnConnected(identifier string)
	OnConnectionFailed(identifier string)
	OnDisconnected(identifier string)
}

// MQTT is the main struct for MQTT messaging
type MQTT struct {
	listeners   map[string]MQTTListener
	isConnected bool
	dialer      net.Dialer
	connection  *net.Conn
}

func NewMQTT() *MQTT {
	mqttInstance := MQTT{
		listeners:   make(map[string]MQTTListener),
		isConnected: false,
		dialer:      net.Dialer{Timeout: defaultConnectTimeout},
		connection:  nil,
	}

	return &mqttInstance
}

// SetConnectTimeout sets the specified conneciton timeout
// which differs from the default time out
func (m *MQTT) SetConnectTimeout(timeout time.Duration) {
	m.dialer.Timeout = timeout
}

// AddListener adds a listener to the mqtt instance
func (m *MQTT) AddListener(listener MQTTListener) {
	if _, ok := m.listeners[listener.Identifier()]; !ok {
		m.listeners[listener.Identifier()] = listener
	}
}

// RemoveListener removes a listener from the mqtt instance
func (m *MQTT) RemoveListener(listener MQTTListener) error {
	if _, ok := m.listeners[listener.Identifier()]; ok {
		delete(m.listeners, listener.Identifier())
		return nil
	}
	return fmt.Errorf("unable to remove listener: %+v", listener.Identifier())
}

// Connect will connect to the given server
func (m *MQTT) Connect(url string, port uint16, useTLS bool) error {
	targetURL := net.JoinHostPort(url, fmt.Sprintf("%d", port))
	ctx, cancel := context.WithTimeout(context.Background(), m.dialer.Timeout)
	defer cancel()

	connection, err := m.dialer.DialContext(ctx, "tcp", targetURL)
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
