package grpcserver

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/protoc/pbs"
	"slot_server/router"
)

type MemeBattleService struct {
	pbs.UnimplementedMemeBattleServiceServer
}

func (rs *MemeBattleService) WayRouterInit() {

}

func (s *MemeBattleService) mustEmbedUnimplementedMemeBattleServiceServer() {
	//TODO implement me
	panic("implement me")
}

//func (s *MemeBattleService) RatTest(ctx context.Context, req *pbs.RatTestReq) (*pbs.NetMessage, error) {
//	fmt.Println("RatTest === ", req)
//	res := "{name:123 ,str:34}"
//	resMarshal, _ := json.Marshal(res)
//	resp := &pbs.NetMessage{
//		ReqHead:   nil,
//		AckHead:   nil,
//		ServiceId: "",
//		MsgId:     0,
//		Content:   resMarshal,
//	}
//	fmt.Println("RatTest succ and return ...")
//	return resp, nil
//}

func (s *MemeBattleService) BeforeExec(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	var (
		comResp = &pbs.NetMessage{
			ReqHead: &pbs.ReqHead{},
			AckHead: &pbs.AckHead{
				Uid:     0,
				Code:    pbs.Code_OK,
				Message: "",
			},
			ServiceId: "",
			MsgId:     0,
			Content:   make([]byte, 0),
		}
	)
	//在执行具体业务前 做一下通用的限制
	//global.GVA_LOG.Infof("BeforeExec:限制协议 %v", req.MsgId)
	//code := logic.GetGameServiceApisAuth(int(req.MsgId))
	//if code != pbs.Code_OK {
	//	comResp.AckHead.Code = code
	//	return comResp, nil
	//}
	return comResp, nil
}

// ComEntranceFunc 通用入口
func (s *MemeBattleService) ComEntranceFunc(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	//global.GVA_LOG.Infof("通用入口 comEntranceFunc :%v", req)

	//根据协议号匹配相应的路由
	// 采用 map 注册的方式
	if value, ok := router.GetHandlersProto(req.MsgId); ok {
		//执行业务前
		comResp, err := s.BeforeExec(ctx, req)
		if err != nil || comResp.AckHead.Code != pbs.Code_OK {
			return comResp, err
		}
		//执行业务接口
		comResp, err = value(ctx, req)
		if err != nil {
			return comResp, err
		}

		//global.GVA_LOG.Infof("通用入口 comEntranceFunc comResp:%v ", comResp)
		return comResp, nil
	} else {
		global.GVA_LOG.Error("comEntranceFunc,处理数据,路由不存在 ", zap.Any("req", req))
		return &pbs.NetMessage{}, errors.New("路由不存在")
	}

}
