package component

import "slot_server/protoc/pbs"

func NewNetMessage(msgId int32) *pbs.NetMessage {
	return &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{},
		AckHead: &pbs.AckHead{
			Uid:     "",
			Code:    pbs.Code_OK,
			Message: "",
		},
		ServiceId: "",
		MsgId:     msgId,
		Content:   make([]byte, 0),
	}
}
