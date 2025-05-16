`# 脚本服务器
`
##  项目编译
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o slot_server_linux
    scp /Users/syn/goProjects/zhuoshuo/nb_game_server/meme_battle/slot_server_linux songyanan@47.97.182.65:/home/songyanan/www/meme_battle
    nohup ./meme_battle_linux  > /dev/null  2>&1 &
