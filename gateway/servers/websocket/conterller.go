package websocket

import (
	"gateway/common"
	"gateway/config"
	"gateway/protoc/pbs"
)

func MemeEntry(message []byte, uid string, msgId int32) (code uint32, msg string, data interface{}) {
	code = common.OK

	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      uid,
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     "",
			Code:    0,
			Message: "",
		},
		ServiceId: config.NatsSlotServer,
		MsgId:     msgId,
		Content:   message,
	}

	//组装
	NastManager.SendMemeJs(&msgReq)

	return code, "", nil
}
