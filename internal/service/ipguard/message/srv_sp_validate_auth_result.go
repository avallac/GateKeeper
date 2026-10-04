package message

const SrvSpValidateAuthResultHeder = "srv_sp_validate_auth_result"

type SrvSpValidateAuthResultBody struct {
	BaseBody
	Result int `json:"result"`
}

func CreateSrvSpValidateAuthResultMessage(result int) ([]byte, error) {
	body := SrvSpValidateAuthResultBody{
		BaseBody: NewBaseBody(),
		Result:   result,
	}
	return Encode(New(SrvSpValidateAuthResultHeder, "mockserver", body))
}
