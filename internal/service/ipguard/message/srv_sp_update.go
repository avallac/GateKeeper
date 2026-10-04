package message

const SrvSpUpdateHeader = "srv_sp_update"

type SrvSpUpdateMessage = Message[SrvSpUpdateBody]

type SrvSpUpdateBody struct {
	BaseBody
	Data []map[string]string `json:"data"`
}
