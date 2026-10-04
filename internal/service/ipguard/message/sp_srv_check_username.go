package message

const SpSrvCheckUsernameHeader = "sp_srv_check_username"

type SpSrvCheckUsernameMessage = Message[SpSrvCheckUsernameBody]

type SpSrvCheckUsernameBody struct {
	BaseBody
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Username string `json:"username"`
	Model    string `json:"model"`
}

func CreateSpSrvCheckUsernameMessage(sender, phone, code, username string) ([]byte, error) {
	body := SpSrvCheckUsernameBody{
		BaseBody: NewBaseBody(),
		Phone:    phone,
		Code:     code,
		Username: username,
		Model:    "sdk_gphone64_arm64",
	}
	return Encode(New(SpSrvCheckUsernameHeader, sender, body))
}
