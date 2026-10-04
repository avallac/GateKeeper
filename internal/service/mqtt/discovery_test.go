package mqtt

import (
	"encoding/json"
	"testing"

	"github.com/avallac/gatekeeper/internal/service/ipguard"
)

var fixtureBoard = ipguard.AppBoard{
	Identity: "00000000-1111-2222-3333-444444444444",
	Name:     "Main Entrance",
}

const fixtureSite = "Test Place"

func TestButtonPayload(t *testing.T) {
	payload := buttonPayload(fixtureBoard, fixtureSite)

	if payload.UniqueID != "00000000-1111-2222-3333-444444444444_open" {
		t.Errorf("UniqueID = %q, want fixture_open suffix", payload.UniqueID)
	}
	if payload.CommandTopic != "gatekeeper/00000000-1111-2222-3333-444444444444/button/set" {
		t.Errorf("CommandTopic = %q", payload.CommandTopic)
	}
	if payload.PayloadPress != "OPEN" {
		t.Errorf("PayloadPress = %q, want OPEN", payload.PayloadPress)
	}
	if len(payload.Device.Identifiers) != 1 || payload.Device.Identifiers[0] != fixtureBoard.Identity {
		t.Errorf("Device.Identifiers = %+v, want [%q]", payload.Device.Identifiers, fixtureBoard.Identity)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var round map[string]any
	if err := json.Unmarshal(data, &round); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if _, ok := round["state_topic"]; ok {
		t.Errorf("button discovery payload must not have a state_topic (stateless)")
	}
}

func TestEventPayload(t *testing.T) {
	payload := eventPayload(fixtureBoard, fixtureSite)

	if payload.StateTopic != "gatekeeper/00000000-1111-2222-3333-444444444444/event" {
		t.Errorf("StateTopic = %q", payload.StateTopic)
	}
	if len(payload.EventTypes) != 1 || payload.EventTypes[0] != "ring" {
		t.Errorf("EventTypes = %+v, want [ring]", payload.EventTypes)
	}
}

func TestCameraPayload(t *testing.T) {
	payload := cameraPayload(fixtureBoard, fixtureSite)

	if payload.Topic != "gatekeeper/00000000-1111-2222-3333-444444444444/camera" {
		t.Errorf("Topic = %q", payload.Topic)
	}
}

func TestDiscoveryTopic(t *testing.T) {
	got := discoveryTopic("button", fixtureBoard.Identity, "open")
	want := "homeassistant/button/00000000-1111-2222-3333-444444444444_open/config"
	if got != want {
		t.Errorf("discoveryTopic() = %q, want %q", got, want)
	}
}
