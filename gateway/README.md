# slot_game_server

##  项目编译
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o  gateway_linux
    scp /Users/syn/goProjects/zhuoshuo/nb_game_server/gateway_sys/gateway_linux songyanan@47.97.182.65:/home/songyanan/www/meme_gate_way
    nohup ./gateway_linux  > /dev/null  2>&1 &

##  swag 文档
    1：每次修改完接口 需要更新swag文档 执行命令
        swag init
    scp /Users/syn/goProjects/zhuoshuo/nb_game_server/gateway_sys/docs/* songyanan@47.97.182.65:/home/songyanan/www/gateway_sys/docs/
