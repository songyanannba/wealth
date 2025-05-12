`# 脚本服务器
`
##  项目编译
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o meme_battle_linux
    scp /Users/syn/goProjects/zhuoshuo/nb_game_server/meme_battle/meme_battle_linux songyanan@47.97.182.65:/home/songyanan/www/meme_battle
    nohup ./meme_battle_linux  > /dev/null  2>&1 &
