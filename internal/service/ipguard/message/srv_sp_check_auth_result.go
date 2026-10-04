package message

const SrvSpCheckAuthResultHeader = "srv_sp_check_auth_result"

type SrvSpCheckAuthResulMessage = Message[SrvSpCheckAuthResultBody]
type SrvSpCheckAuthResultBody struct {
	BaseBody
	Result int `json:"result"`
}

func CreateSrvSpCheckAuthResultMessage(result int) ([]byte, error) {
	body := SrvSpCheckAuthResultBody{
		BaseBody: NewBaseBody(),
		Result:   result,
	}
	return Encode(New(SrvSpCheckAuthResultHeader, "fake-server", body))
}
