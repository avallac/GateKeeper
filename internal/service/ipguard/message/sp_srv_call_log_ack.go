package message

const SpSrvCallLogAckHeader = "sp_srv_call_log_ack"

type SpSrvCallLogAckMessage = Message[SpSrvCallLogAckBody]

type SpSrvCallLogAckBody struct {
	BaseBody
	Token string `json:"token"`
}

func CreateSpSrvCallLogAckMessage(sender, token string) ([]byte, error) {
	body := SpSrvCallLogAckBody{
		BaseBody: NewBaseBody(),
		Token:    token,
	}
	return Encode(New(SpSrvCallLogAckHeader, sender, body))
}
