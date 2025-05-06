.PHONY: run


# 生成游戏rat proto文件
pbs-rat:
	@protoc --proto_path=./ --go_out=./pbs --go_opt=paths=source_relative --go-grpc_out=./pbs --go-grpc_opt=paths=source_relative *.proto

