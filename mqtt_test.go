// Package mqtt is the main encompasing mqtt package
package mqtt

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	testURL  = "test.mosquitto.org"
	testPort = 1883
)

type MQTTTest struct {
	suite.Suite
	instance *MQTT
}

func (m *MQTTTest) SetupTest() {
	m.instance = &MQTT{isConnected: false}
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

func TestMQTTSuite(t *testing.T) {
	suite.Run(t, new(MQTTTest))
}
