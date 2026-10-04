package message

const SpSrvAskKeyHeader = "sp_srv_ask_key"

type SpSrvAskKeyMessage = Message[SpSrvAskKey]

type SpSrvAskKey struct {
	BaseBody
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func CreateSpSrvAskKeyMessage(sender string, phone string, code string) ([]byte, error) {
	body := SpSrvAskKey{
		BaseBody: NewBaseBody(),
		Phone:    phone,
		Code:     code,
	}
	return Encode(New(SpSrvAskKeyHeader, sender, body))
}
