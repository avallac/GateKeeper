package message

const BcSpRelaysHeader = "bc_sp_relays"

type BcSpRelaysMessage = Message[BcSpRelaysBody]

type BcSpRelaysBody struct {
	BaseBody
	Sites map[string]BcSpRelaysSite `json:"sites"`
}

type BcSpRelaysSite struct {
	Name   string            `json:"name"`
	Boards []BcSpRelaysBoard `json:"boards"`
}

type BcSpRelaysBoard struct {
	Identity                string                     `json:"identity"`
	Name                    string                     `json:"name"`
	DeviceID                string                     `json:"device_id"`
	DeviceType              int                        `json:"device_type"`
	RelayActive             []bool                     `json:"relay_active"`
	RelayAccessMobile       []bool                     `json:"relay_access_mobile"`
	RelayName               []string                   `json:"relay_name"`
	RelayUsage              []int                      `json:"relay_usage"`
	RelayOpenMode           []int                      `json:"relay_open_mode"`
	RelayAccessMobileV2     []BcSpRelaysAccessMobileV2 `json:"relay_access_mobile_v2"`
	Peripherals             []interface{}              `json:"peripherals"`
	Title                   string                     `json:"title"`
	Key                     string                     `json:"key"`
	WaitConfirmBadgeDigital bool                       `json:"waitConfirmBadgeDigital"`
}

type BcSpRelaysAccessMobileV2 struct {
	Type int `json:"type"`
	Mode int `json:"mode"`
}

func CreateBcSpRelaysMessage() ([]byte, error) {
	body := BcSpRelaysBody{
		BaseBody: NewBaseBody(),
		Sites: map[string]BcSpRelaysSite{
			"100500": {
				Name: "Test Place",
				Boards: []BcSpRelaysBoard{
					{
						Identity:          "00000000-1111-2222-3333-444444444444",
						Name:              "Main Entrance",
						DeviceID:          "00:11:22:33:44:55",
						DeviceType:        11,
						RelayActive:       []bool{true, false, false},
						RelayAccessMobile: []bool{true, false, false},
						RelayName:         []string{"Door", "Relay 2", "Relay 3"},
						RelayUsage:        []int{2, 2, 2},
						RelayOpenMode:     []int{0, 0, 0},
						RelayAccessMobileV2: []BcSpRelaysAccessMobileV2{
							{Type: 0, Mode: 0},
							{Type: -1, Mode: -1},
							{Type: -1, Mode: -1},
						},
						Peripherals:             []interface{}{},
						Title:                   "",
						Key:                     "",
						WaitConfirmBadgeDigital: false,
					},
				},
			},
		},
	}
	return Encode(New(BcSpRelaysHeader, "mockserver", body))
}
