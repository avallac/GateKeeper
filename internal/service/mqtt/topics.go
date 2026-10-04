package mqtt

import "fmt"

const bridgeAvailabilityTopic = "gatekeeper/bridge/availability"

const payloadPress = "OPEN"

func discoveryTopic(component, boardID, suffix string) string {
	return fmt.Sprintf("homeassistant/%s/%s_%s/config", component, boardID, suffix)
}

func commandTopic(boardID string) string {
	return fmt.Sprintf("gatekeeper/%s/button/set", boardID)
}

func eventTopic(boardID string) string {
	return fmt.Sprintf("gatekeeper/%s/event", boardID)
}

func cameraTopic(boardID string) string {
	return fmt.Sprintf("gatekeeper/%s/camera", boardID)
}
