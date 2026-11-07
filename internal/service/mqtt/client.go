package mqtt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/avallac/gatekeeper/internal/service"
	"github.com/avallac/gatekeeper/internal/service/ipguard"
	"github.com/avallac/gatekeeper/internal/service/ipguard/message"
)

type Client struct {
	paho pahomqtt.Client
	app  *ipguard.Client
}

func NewClient(appConfig *service.AppConfig, app *ipguard.Client) (*Client, error) {
	if appConfig.MQTTHost == "" {
		return nil, fmt.Errorf("no MQTT broker configured")
	}

	scheme := "tcp"
	if appConfig.MQTTSSL {
		scheme = "ssl"
	}
	port := appConfig.MQTTPort
	if port == "" {
		port = "1883"
	}
	broker := fmt.Sprintf("%s://%s:%s", scheme, appConfig.MQTTHost, port)

	client := &Client{app: app}

	app.OnRelaysUpdated(func(_ []ipguard.AppSite) {
		client.publishAllDiscovery()
		client.subscribeCommands()
	})
	app.OnCallLog(func(boardID string, _ ipguard.AppLog, body message.BrdSpCallLogBody) {
		client.publishPhotoEvent(boardID, body)
	})
	app.OnVideoCall(func(boardID string) {
		client.publishRingEvent(boardID)
	})

	opts := pahomqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("gatekeeper-"+app.UserConfig.Sender).
		SetUsername(appConfig.MQTTUsername).
		SetPassword(appConfig.MQTTPassword).
		SetAutoReconnect(true).
		SetWill(bridgeAvailabilityTopic, "offline", 1, true)
	opts.SetOnConnectHandler(func(_ pahomqtt.Client) {
		log.Println("[mqtt] connected")
		client.publish(bridgeAvailabilityTopic, 1, true, "online")
		client.publishAllDiscovery()
		client.subscribeCommands()
	})

	client.paho = pahomqtt.NewClient(opts)
	token := client.paho.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) publishAllDiscovery() {
	for _, site := range c.app.Sites {
		for _, board := range site.Boards {
			c.publishDiscovery("button", board.Identity, "open", buttonPayload(board, site.Name))
			c.publishDiscovery("event", board.Identity, "ring", ringEventPayload(board, site.Name))
			c.publishDiscovery("event", board.Identity, "photo", photoEventPayload(board, site.Name))
			c.publishDiscovery("camera", board.Identity, "snapshot", cameraPayload(board, site.Name))
		}
	}
}

func (c *Client) publishDiscovery(component, boardID, suffix string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[mqtt] discovery marshal failed: %v", err)
		return
	}
	c.publish(discoveryTopic(component, boardID, suffix), 1, true, data)
}

func (c *Client) subscribeCommands() {
	for _, site := range c.app.Sites {
		for _, board := range site.Boards {
			id := board.Identity
			token := c.paho.Subscribe(commandTopic(id), 1, func(_ pahomqtt.Client, msg pahomqtt.Message) {
				if string(msg.Payload()) == payloadPress {
					c.app.OpenDoor(id)
				}
			})
			token.Wait()
			if err := token.Error(); err != nil {
				log.Printf("[mqtt] subscribe %s failed: %v", commandTopic(id), err)
			}
		}
	}
}

func (c *Client) publishPhotoEvent(boardID string, body message.BrdSpCallLogBody) {
	if body.PhotoData != "" {
		raw, err := base64.StdEncoding.DecodeString(body.PhotoData)
		if err != nil {
			log.Printf("[mqtt] failed to decode call photo: %v", err)
		} else {
			c.publish(cameraTopic(boardID), 0, false, raw)
		}
	}
	payload, _ := json.Marshal(map[string]any{
		"event_type":     "photo",
		"responder_name": body.ResponderName,
		"ts":             body.Timestring,
	})
	c.publish(photoTopic(boardID), 0, false, payload)
}

func (c *Client) publishRingEvent(boardID string) {
	payload, _ := json.Marshal(map[string]any{
		"event_type": "ring",
	})
	c.publish(ringTopic(boardID), 0, false, payload)
}

func (c *Client) publish(topic string, qos byte, retained bool, payload any) {
	token := c.paho.Publish(topic, qos, retained, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("[mqtt] publish %s failed: %v", topic, err)
	}
}
