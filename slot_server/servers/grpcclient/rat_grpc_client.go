package grpcclient

//
//var (
//	clientInstance *GRPCClient
//	once           sync.Once
//)
//
//// GRPCClient 封装 gRPC 客户端和连接
//type GRPCClient struct {
//	conn *grpc.ClientConn
//	//RatServiceClient rat.RatServiceClient // 替换为你的服务客户端
//	ServiceClient pbs.RatServiceClient // 替换为你的服务客户端
//}
//
//// InitClient 初始化全局 gRPC 客户端
//func InitClient(address string) {
//	once.Do(func() {
//		conn, err := grpc.Dial(address, grpc.WithInsecure())
//		if err != nil {
//			log.Fatalf("failed to connect: %v", err)
//		}
//		clientInstance = &GRPCClient{
//			conn:          conn,
//			ServiceClient: pbs.NewRatServiceClient(conn), // 替换为你的服务客户端
//		}
//	})
//}
//
//// GetClient 获取全局 gRPC 客户端实例
//func GetClient() *GRPCClient {
//	return clientInstance
//}
//
//// Close 关闭全局 gRPC 客户端连接
//func Close() {
//	if clientInstance != nil && clientInstance.conn != nil {
//		clientInstance.conn.Close()
//	}
//}
//
//// CallMethod 封装 gRPC 方法调用
//func (g *GRPCClient) CallMethod(req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	// 创建上下文并设置超时
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	entranceFunc, err := g.ServiceClient.ComEntranceFunc(ctx, req)
//
//	global.GVA_LOG.Infof("CallMethod entranceFunc:%v, err%v", entranceFunc, err)
//	return entranceFunc, err
//}
