package message

const SpSrvCheckAuthHeader = "sp_srv_check_auth"

type SpSrvCheckAuthBody struct {
	BaseBody
	Phone string `json:"phone"`
}

func CreateSpSrvCheckAuthMessage(sender, phone string) ([]byte, error) {
	body := SpSrvCheckAuthBody{
		BaseBody: NewBaseBody(),
		Phone:    phone,
	}
	return Encode(New(SpSrvCheckAuthHeader, sender, body))
}
