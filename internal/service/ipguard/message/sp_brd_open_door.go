package message

const SpBrdOpenDoorHeader = "sp_brd_open_door"

type SpBrdOpenDoorMessage = Message[SpBrdOpenDoorBody]

type SpBrdOpenDoorBody struct {
	BaseBody
	Door int `json:"door"`
}

func CreateOpenDoorMessage(sender string) ([]byte, error) {
	body := SpBrdOpenDoorBody{
		BaseBody: NewBaseBody(),
		Door:     1,
	}
	return Encode(New(SpBrdOpenDoorHeader, sender, body))
}
