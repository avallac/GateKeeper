package message

const SpSrvUpdateAckHeader = "sp_srv_update_ack"

type SpSrvUpdateAckMessage = Message[SpSrvUpdateAckBody]

type SpSrvUpdateAckBody struct {
	BaseBody
}

func CreateSrvUpdateAckMessage(sender, token string) ([]byte, error) {
	body := SpSrvUpdateAckBody{
		BaseBody: BaseBody{Token: token},
	}
	return Encode(New(SpSrvUpdateAckHeader, sender, body))
}
