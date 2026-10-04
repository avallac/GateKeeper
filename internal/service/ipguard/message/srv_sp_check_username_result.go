package message

const SrvSpCheckUsernameResultHeader = "srv_sp_check_username_result"

type SrvSpCheckUsernameResultMessage = Message[SrvSpCheckUsernameResultBody]
type SrvSpCheckUsernameResultBody struct {
	BaseBody
	Result int `json:"result"`
}

func CreateSrvSpCheckUsernameResultMessage(result int) ([]byte, error) {
	body := SrvSpCheckAuthResultBody{
		BaseBody: NewBaseBody(),
		Result:   result,
	}
	return Encode(New(SrvSpCheckUsernameResultHeader, "fake-server", body))
}
