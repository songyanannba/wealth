.PHONY: run
run:
	go run main.go

# 生成通用proto文件
#pbs-common:
#	@protoc --go_out=../ ./proto/common/*.proto
#
## 生成游戏proto文件
#pbs-game:
#	@protoc --go_out=../ ./proto/game/*.proto
#
## 生成所有proto
#pbs-all:
#	@protoc --go_out=../ ./proto/common/*.proto ./proto/game/*.proto

#构建服务
build-linux-serve:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gate_way_linux


dev_scp-build-linux-serve:
    scp /Users/syn/goProjects/zhuoshuo/nb_game_server/gate_way/gate_way_linux songyanan@47.97.182.65:/home/songyanan/www/gate_way

#上传swagger 文档
pro_scp-build-linux-serve:
    scp /Users/syn/goProjects/zhuoshuo/nb_game_server/gate_way/docs/* songyanan@47.97.182.65:/home/songyanan/www/gate_way/docs/


#生产
pro_scp-build-linux-serve:
    ##scp /Users/syn/goProjects/zhuoshuo/Battle-Royale/dataosha-server-linux songyanan@121.199.26.144:/home/songyanan/www/da_tao_sha



