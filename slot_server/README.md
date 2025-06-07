`# 脚本服务器
`
##  项目编译
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o slot_server_linux
    nohup ./slot_server_linux  > /dev/null  2>&1 &
