package message

const SrvSpAskKeyResultHeader = "srv_sp_ask_key_result"

type SrvSpAskKeyResultMessage = Message[SrvSpAskKeyResultBody]

type SrvSpAskKeyResultBody struct {
	BaseBody
	Result int    `json:"result"`
	Key    string `json:"key"`
}

func CreateSrvSpAskKeyResultMessage(status int, key string) ([]byte, error) {
	body := SrvSpAskKeyResultBody{
		BaseBody: NewBaseBody(),
		Result:   status,
		Key:      key,
	}
	return Encode(New(SrvSpAskKeyResultHeader, "fake-server", body))
}
