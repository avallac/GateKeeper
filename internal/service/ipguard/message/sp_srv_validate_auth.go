package message

const SpSrvValidateAuthHeader = "sp_srv_validate_auth"

type SpSrvValidateAuthMessage = Message[SpSrvValidateAuthBody]

type SpSrvValidateAuthBody struct {
	BaseBody
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func CreateSpSrvValidateAuthMessage(sender string, phone string, code string) ([]byte, error) {
	body := SpSrvValidateAuthBody{
		BaseBody: NewBaseBody(),
		Phone:    phone,
		Code:     code,
	}
	return Encode(New(SpSrvValidateAuthHeader, sender, body))
}
