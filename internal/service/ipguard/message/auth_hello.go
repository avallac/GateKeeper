package message

type AuthHelloMessage = Message[AuthHelloBody]

type AuthHelloBody struct {
	BaseBody
	Type    int               `json:"type"`
	Version string            `json:"version"`
	Data    AuthHelloBodyData `json:"data"`
}

type AuthHelloBodyData struct {
	FCMToken   string     `json:"fcm_token"`
	WithoutSIP bool       `json:"without_sip"`
	OSVersion  string     `json:"os_version"`
	Language   string     `json:"language"`
	Timezone   int        `json:"timezone"`
	Weektypes  [][]string `json:"weektypes"`
}

func CreateAuthHelloMessage(sender, version string) ([]byte, error) {
	weekTypes := make([][]string, 7)
	for i := 0; i < 7; i++ {
		weekTypes[i] = []string{}
	}

	body := AuthHelloBody{
		BaseBody: NewBaseBody(),
		Type:     0,
		Version:  version,
		Data: AuthHelloBodyData{
			FCMToken:   "",
			WithoutSIP: true,
			OSVersion:  "16",
			Language:   "en_US",
			Timezone:   3600,
			Weektypes:  weekTypes,
		},
	}

	return Encode(New(HelloHeader, sender, body))
}
