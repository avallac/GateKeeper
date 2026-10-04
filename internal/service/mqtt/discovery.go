package mqtt

import "github.com/avallac/gatekeeper/internal/service/ipguard"

type device struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model,omitempty"`
}

type buttonDiscovery struct {
	Name              string `json:"name"`
	UniqueID          string `json:"unique_id"`
	CommandTopic      string `json:"command_topic"`
	PayloadPress      string `json:"payload_press"`
	AvailabilityTopic string `json:"availability_topic"`
	Device            device `json:"device"`
}

type eventDiscovery struct {
	Name              string   `json:"name"`
	UniqueID          string   `json:"unique_id"`
	StateTopic        string   `json:"state_topic"`
	EventTypes        []string `json:"event_types"`
	AvailabilityTopic string   `json:"availability_topic"`
	Device            device   `json:"device"`
}

type cameraDiscovery struct {
	Name              string `json:"name"`
	UniqueID          string `json:"unique_id"`
	Topic             string `json:"topic"`
	AvailabilityTopic string `json:"availability_topic"`
	Device            device `json:"device"`
}

func buildDevice(board ipguard.AppBoard, siteName string) device {
	return device{
		Identifiers:  []string{board.Identity},
		Name:         board.Name,
		Manufacturer: "IPGuard",
		Model:        siteName,
	}
}

func buttonPayload(board ipguard.AppBoard, siteName string) buttonDiscovery {
	return buttonDiscovery{
		Name:              board.Name,
		UniqueID:          board.Identity + "_open",
		CommandTopic:      commandTopic(board.Identity),
		PayloadPress:      payloadPress,
		AvailabilityTopic: bridgeAvailabilityTopic,
		Device:            buildDevice(board, siteName),
	}
}

func eventPayload(board ipguard.AppBoard, siteName string) eventDiscovery {
	return eventDiscovery{
		Name:              board.Name + " Call",
		UniqueID:          board.Identity + "_ring",
		StateTopic:        eventTopic(board.Identity),
		EventTypes:        []string{"ring"},
		AvailabilityTopic: bridgeAvailabilityTopic,
		Device:            buildDevice(board, siteName),
	}
}

func cameraPayload(board ipguard.AppBoard, siteName string) cameraDiscovery {
	return cameraDiscovery{
		Name:              board.Name + " Snapshot",
		UniqueID:          board.Identity + "_snapshot",
		Topic:             cameraTopic(board.Identity),
		AvailabilityTopic: bridgeAvailabilityTopic,
		Device:            buildDevice(board, siteName),
	}
}
