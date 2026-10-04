package message

const SpSrvOverrideUsernameHeader = "sp_srv_override_username"

type SpSrvOverrideUsernameMessage = Message[SpSrvOverrideUsernameBody]

type SpSrvOverrideUsernameBody struct {
	BaseBody
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Username string `json:"username"`
	Model    string `json:"model"`
}

func CreateSpSrvOverrideUsernameBodyMessage(sender, phone, code, username string) ([]byte, error) {
	body := SpSrvOverrideUsernameBody{
		BaseBody: NewBaseBody(),
		Phone:    phone,
		Code:     code,
		Username: username,
		Model:    "sdk_gphone64_arm64",
	}
	return Encode(New(SpSrvOverrideUsernameHeader, sender, body))
}
