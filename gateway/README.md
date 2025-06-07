# slot_game_server

##  项目编译
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o  gateway_linux
    nohup ./gateway_linux  > /dev/null  2>&1 &


