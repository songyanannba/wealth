
## 下面是短链接注册登录获取用户信息 

# 注册
```
curl --location --request POST '47.97.201.179:8000/register' \
--header 'Content-Type: application/json' \
--data-raw '{
"user_name": "11111",
"pass_word": "111"
}'
```
```json
{
  "code": 200,
  "msg": "Success",
  "data": {}
}
```

# 登录
```
curl --location --request POST '47.97.201.179:8000/login' \
--header 'Content-Type: application/json' \
--data-raw '{
"user_name": "1111",
"pass_word": "111"
}'
```
```json

{
  "code": 200,
  "msg": "Success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY3Nzg5ODksInN1YiI6IiIsInVzZXJfaWQiOiI0NjVhNTE2ZC00OWQ2LTQxNDMtODNkZS0wY2M4NjU1MjNlMGIifQ.LUQfjWQHf3nHR85PqkOAWHuC9uHO_1dsYmzadatmvGw"
  }
}
```


# 短链接获取用户信息 (请求头需要添加 gw-token 字段)
```
curl --location --request GET '47.97.201.179:8000/game_way/user_info' \
--header 'gw-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY3Nzg5ODksInN1YiI6IiIsInVzZXJfaWQiOiI0NjVhNTE2ZC00OWQ2LTQxNDMtODNkZS0wY2M4NjU1MjNlMGIifQ.LUQfjWQHf3nHR85PqkOAWHuC9uHO_1dsYmzadatmvGw'

```
```json
{
"code": 200,
"msg": "Success",
"data": {
"user_info": {
"user_name": "11111",
"amount": 0,
"city": "china"
}
}
}
```


### 完成注册后 需要使用websocket 长链接 
# ws://47.97.201.179:8081/gate_way （链接服务器的时候需要在请求头加 gw-token 字段，这个字段是短链接登录后获取到的）


# 下面的请求和返回结构具体数据都在   message netMessage {} ； 里面的 content ； netMessage是统一的请求和返回结构

# 5月16日 -----

# 主动发送消息
##  协议 10004   1:获取当前局动物派对的信息 ; 请求结构  message currAPInfoReq {} ； 返回结构体查看 :  message currAPInfoAck {}
##  协议 100011  2:押注                  ; 请求结构：  message userBetReq {} ； 返回结构体查看 :    message userBetAck {}
##  协议 10009   3:获取所有的在线用户

# 被动收到消息
## 协议 10006： //押注停止后 主动下发最外圈的动物排序，第一个排在最上面 位置0开始  ； 结构体查看 : message  animalSortMsg {}
## 协议 10007： //当前期是否赢钱的通知 ;  结构体查看 : message  currPeriodUserWinMsg {}
## 协议 10008： //当前轮结束 下一轮开始发送颜色列表 ; 结构体查看 : message  colorSortMsg {}

# 5月16日 -----