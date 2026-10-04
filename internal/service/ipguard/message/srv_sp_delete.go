package message

const SrvSpDeleteHeader = "srv_sp_delete"

type SrvSpDeleteMessage = Message[SrvSpDeleteBody]

type SrvSpDeleteBody struct {
	BaseBody
}
