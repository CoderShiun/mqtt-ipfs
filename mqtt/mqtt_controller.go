package mqtt

import (
	"fmt"
	"sync"
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	Host     = "127.0.0.1:1883"//"192.168.1.101:8000"
	UserName = "mxc"
	Password = "mxc11908"
)

type Client struct {
	nativeClient  gomqtt.Client
	clientOptions *gomqtt.ClientOptions
	locker        *sync.Mutex
	// 消息收到之后处理函数
	observer func(c *Client, msg *Message)
}

type Message struct {
	ClientID string `json:"clientId"`
	Type     string `json:"type"`
	Data     string `json:"data,omitempty"`
	Time     int64  `json:"time"`
}

func NewClient(clientId string) *Client {
	clientOptions := gomqtt.NewClientOptions().
		AddBroker(Host).
		SetUsername(UserName).
		SetPassword(Password).
		SetClientID(clientId).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetKeepAlive(120 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetWriteTimeout(10 * time.Second).
		SetOnConnectHandler(func(client gomqtt.Client) {
			// 连接被建立后的回调函数
			fmt.Println("Mqtt is connected!", "clientId", clientId)
		}).
		SetConnectionLostHandler(func(client gomqtt.Client, err error) {
			// 连接被关闭后的回调函数
			fmt.Println("Mqtt is disconnected!", "clientId", clientId, "reason", err.Error())
		})

	nativeClient := gomqtt.NewClient(clientOptions)

	return &Client{
		nativeClient:  nativeClient,
		clientOptions: clientOptions,
		locker:        &sync.Mutex{},
	}
}

func (client *Client) GetClientID() string {
	return client.clientOptions.ClientID
}

func (client *Client) Connect() error {
	return client.ensureConnected()
}

func (client *Client) ensureConnected() error {
	if !client.nativeClient.IsConnected() {
		client.locker.Lock()
		defer client.locker.Unlock()
		if !client.nativeClient.IsConnected() {
			if token := client.nativeClient.Connect(); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
	}
	return nil
}