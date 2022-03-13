// Package mqtt is the main encompasing mqtt package
package mqtt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	testURL     = "test.mosquitto.org"
	testPort    = 1883
	oneListener = 1
	noListeners = 0
)

type MQTTTest struct {
	suite.Suite
	instance *MQTT
}

type DefaultListener struct {
	name string
}

func (d *DefaultListener) Identifier() string {
	return d.name
}

func (d *DefaultListener) OnConnected(identifier string) {
	fmt.Printf("connected: %s", identifier)
}

func (d *DefaultListener) OnDisconnected(identifier string) {
	fmt.Printf("disconnected: %s", identifier)
}

func (d *DefaultListener) OnConnectionFailed(identifier string) {
	fmt.Printf("connection failed: %s", identifier)
}

func (m *MQTTTest) SetupTest() {
	m.instance = NewMQTT()
}

func (m *MQTTTest) TestConnect() {
	err := m.instance.Connect(testURL, testPort, false)
	m.Assert().Nil(err)
	m.Assert().True(m.instance.isConnected)
}

func (m *MQTTTest) TestDisconnect() {
	err := m.instance.Connect(testURL, testPort, false)
	m.Assert().Nil(err)

	err = m.instance.Disconnect()
	m.Assert().Nil(err)

	m.Assert().Equal(false, m.instance.isConnected)
}

func (m *MQTTTest) TestListeners() {
	defaultListener := DefaultListener{name: "default"}
	m.instance.AddListener(&defaultListener)
	m.Assert().Equal(oneListener, len(m.instance.listeners))

	err := m.instance.RemoveListener(&defaultListener)
	m.Assert().Nil(err)
	m.Assert().Equal(noListeners, len(m.instance.listeners))
}

func (m *MQTTTest) TestCallbacks() {
	defaultListener := DefaultListener{name: "default"}
	m.instance.AddListener(&defaultListener)

	err := m.instance.Connect(testURL, testPort, false)
	m.Assert().Nil(err)
}

func TestMQTTSuite(t *testing.T) {
	suite.Run(t, new(MQTTTest))
}
