package message

const PsSpSdpOfferHeader = "ps_sp_sdp_offer"

type PsSpSdpOfferMessage = Message[PsSpSdpOfferBody]

type PsSpSdpOfferBody struct {
	BaseBody
	RemoteIdentity string `json:"remoteIdentity"`
	SDP            []byte `json:"sdp"`
}
