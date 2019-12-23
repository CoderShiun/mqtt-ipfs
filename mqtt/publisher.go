package mqtt

import (
	"errors"
	"time"
)

// retained: should keep message or not
func (client *Client) Publish(topic string, qos byte, retained bool, data []byte) error {
	if err := client.ensureConnected(); err != nil {
		return err
	}

	token := client.nativeClient.Publish(topic, qos, retained, data)
	if err := token.Error(); err != nil {
		return err
	}

	// return false is the timeout occurred
	if !token.WaitTimeout(time.Second * 10) {
		return errors.New("mqtt publish wait timeout")
	}

	return nil
}