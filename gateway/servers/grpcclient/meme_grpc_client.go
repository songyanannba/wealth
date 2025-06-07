package grpcclient

//
//// MebGRPCClient 表情包大作战
//type MebGRPCClient struct {
//	conn            *grpc.ClientConn
//	MTServiceClient pbs.MemeBattleServiceClient // 替换为你的服务客户端
//}
//
//var (
//	mebClientInstance *MebGRPCClient
//	mTowerOnce        sync.Once
//)
//
//// InitMebClient 初始化全局 gRPC 客户端
//func InitMebClient(address string) {
//	mTowerOnce.Do(func() {
//		conn, err := grpc.Dial(address, grpc.WithInsecure())
//		if err != nil {
//			log.Fatalf("failed to connect: %v", err)
//		}
//		mebClientInstance = &MebGRPCClient{
//			conn:            conn,
//			MTServiceClient: pbs.NewMemeBattleServiceClient(conn), // 替换为你的服务客户端
//		}
//	})
//}
//
//// GetMebClient 获取全局 gRPC客户端实例
//func GetMebClient() *MebGRPCClient {
//	return mebClientInstance
//}
//
//// MebClose 关闭全局 gRPC 客户端连接
//func MebClose() {
//	if mebClientInstance != nil && mebClientInstance.conn != nil {
//		err := mebClientInstance.conn.Close()
//		if err != nil {
//			global.GVA_LOG.Error("failed to close connection: %v", zap.Error(err))
//		}
//	}
//}
//
//// CallMebMethod 封装 gRPC方法调用
//func (g *MebGRPCClient) CallMebMethod(req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	// 创建上下文并设置超时
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	entranceFunc, err := g.MTServiceClient.ComEntranceFunc(ctx, req)
//	global.GVA_LOG.Infof("CallMethod entranceFunc:%v, err%v", entranceFunc, err)
//	return entranceFunc, err
//}
