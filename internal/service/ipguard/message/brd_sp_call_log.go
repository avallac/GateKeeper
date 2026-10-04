package message

const BrdSpCallLogHeader = "brd_sp_call_log"

type BrdSpCallLogMessage = Message[BrdSpCallLogBody]

type BrdSpCallLogBody struct {
	BaseBody
	Timestamp      int64  `json:"timestamp"`
	Timestring     string `json:"timestring"`
	ResponderType  int    `json:"responder_type"`
	ResponderID    string `json:"responder_id"`
	ResponderName  string `json:"responder_name"`
	OpenStatus     bool   `json:"open_status"`
	OpenAux1Status bool   `json:"open_aux1_status"`
	OpenAux2Status bool   `json:"open_aux2_status"`
	PhotoDisabled  bool   `json:"photo_disabled"`
	PhotoData      string `json:"photo_data"`
}
