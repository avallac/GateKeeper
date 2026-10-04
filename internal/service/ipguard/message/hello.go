package message

const HelloHeader = "hello"

type HelloMessage = Message[HelloBody]

type HelloBody struct {
	BaseBody
	Type    int           `json:"type"`
	Version string        `json:"version"`
	Data    HelloBodyData `json:"data"`
}

type HelloBodyData struct {
	Username   string     `json:"username"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	Akey       string     `json:"akey"`
	Code       string     `json:"code"`
	FCMToken   string     `json:"fcm_token"`
	WithoutSIP bool       `json:"without_sip"`
	OSVersion  string     `json:"os_version"`
	Language   string     `json:"language"`
	Timezone   int        `json:"timezone"`
	Weektypes  [][]string `json:"weektypes"`
}

func CreateHelloMessage(sender, version, name, phone, code string) ([]byte, error) {
	weekTypes := make([][]string, 7)
	for i := 0; i < 7; i++ {
		weekTypes[i] = []string{}
	}

	body := HelloBody{
		BaseBody: NewBaseBody(),
		Type:     0,
		Version:  version,
		Data: HelloBodyData{
			Username:   name,
			Phone:      phone,
			Email:      "",
			Akey:       "",
			Code:       code,
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
