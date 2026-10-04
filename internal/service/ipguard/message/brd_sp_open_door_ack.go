package message

const BrdSpOpenDoorAckHeader = "brd_sp_open_door_ack"

type BrdSpOpenDoorAckMessage = Message[BrdSpOpenDoorAckBody]

type BrdSpOpenDoorAckBody struct {
	BaseBody
	Door       int  `json:"door"`
	ActiveCall bool `json:"active_call"`
}

func CreateBrdSpOpenDoorAckMessage(sender string) ([]byte, error) {
	body := BrdSpOpenDoorAckBody{
		BaseBody:   NewBaseBody(),
		Door:       1,
		ActiveCall: false,
	}
	return Encode(New(BrdSpOpenDoorAckHeader, sender, body))
}
