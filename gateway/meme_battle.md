# req 完整格式
```json
  {
  "seq": "登陆",
  "cmd": "login",
  "data": {
    "user_id": "20210",
    "token" : "",
    "app_id": 10 ,
    "nickname": "20210"
  }
}
```

# resp 完整格式 
```json
{
  "seq": "1569080188418-747717",
  "cmd": "login",
  "response": {
    "code": 1,
    "codeMsg": "Ok",
    "data": {
      "target": "",
      "type": "text",
      "msg": "hello",
      "from": "马超"
    }
  }
}
```


# 建立连接  ws://127.0.0.1:8199/gate_way
    下面是路由请求例子
    登陆
    ws.send('{"seq":"2323","cmd":"login","data":{"user_id":"11","app_id":10}}');
    心跳:
    ws.send('{"seq":"2324","cmd":"heartbeat","data":{"user_id":""}}');


    完整格式参考 登陆
    {
        "seq": "1",     //时间戳  暂时不需要定义 传时间戳就行
        "cmd": "login", //方法
        "data": {
            "user_id": "6966",
            "app_id": 10 // 10 代表表情包大作战；
        }
    }



```json
```

```json
```

```json
```


# resp 下发的广播
    1001 //创建房间成功的广播
    1015 //邀请好友的广播（被邀请人会收到）
    1003 //加入房间的广播（mebJoinRoom ，此接口就是正常的加入，断线从连不要调用这个，收到邀请通知调用这个，分享链接调用这个）
    1005 //离开房间的广播 （当房间就一个人 离开房间等同于解散）
    1013 //踢人广播
    1007 //获取用户当前游戏状态的广播 （用户进入游戏后 需要先调用此接口，如果存在房间编号，说明已经加入过，这个时候需要走断线从连逻辑，请求重新加入房间接口：mebReJoinRoom ，注意区分 mebJoinRoom加入接口 ）
    1009 //重新加入房间的广播
    1017 //开始对局游戏
    1019 //加载 （服务端用来确认用户收到1017后 是否进入房间）
    1021 // 问题广播
    1022 // 发手牌广播
    1026 // 出牌成功成功广播
    1027 // 重新随牌成功广播
    1028 // 本轮所有用户出的牌 （收到此消息 进入点赞页面）
    1029 // 点赞
    1030 // 本局结束 计算本局最终结果
    1031 // 游戏结束
    1032 // 匹配
    1033 // 取消匹配 mebCancel
    1034 // 匹配成功 并开始
    1035 // 就绪
    1036 // 图鉴列表
    1037 // 拆包
    1038 // 版本列表
    1055 // 取消准备


    

### 除了登陆和心跳，其他的游戏相关接口，有请求和返回和广播，3部分组成
### 因为是异步处理，需要知道，返回成功，不一定处理成功，只有在返回成功并且收到成功的广播， 才算是处理成功。如果异步处理逻辑失败，会给用户广播失败原因



# === 1 登陆  login
```json
  {
  "seq": "登陆",
  "cmd": "login",
  "data": {
    "user_id": "20210",
    "token" : "",
    "app_id": 10 ,
    "nickname": "20210"
  }
}
```

```json
{
    "seq": "登陆",
    "cmd": "login",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDE1MDQ1NjEsInN1YiI6IiIsInVzZXJfaWQiOiIyMDIxMCJ9.2xdNsd5SKMSXi8eB1-zLNdpef6LBTKQAuA1XfB1-vxg",
            "user_id": "20210",
            "nickname": "20210"
        }
    }
}
```

# ===2 用户心跳  login 创建用户长链接 (登陆后如果不退出 就每隔5秒发一次)
    -- heartbeat
    req
        {
            "user_id": 1  string //用户ID
        }

    resp data
        {

        }

# === 创建房间 

[//]: # (type CreateRoomReq struct {)
[//]: # (UserID       string `json:"user_id"`)
[//]: # (RoomType     int    `json:"room_type"`      // 1:好友约战)
[//]: # (UserNumLimit int    `json:"user_num_limit"` //用户人数限制 2人场 3 人场 4人场)
[//]: # (RoomTurnNum  int    `json:"room_turn_num"`  //房间 回合数 3/5/7)
[//]: # (})

```json
{
  "seq": "创建房间",
  "cmd": "mebCreateRoom",
  "service_id": "meme_battle",
  "data": {
    "room_type": 1,
    "user_num_limit": 2,
    "room_turn_num": 1
  }
}
```

```json
{
    "seq": "创建房间",
    "cmd": "mebCreateRoom",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": null
    }
}
```

## 房间状态参考 (status ): 房间状态: 1=开放中,2=已满员,3=已解散,4=进行中,5=已结束 6=异常房间 7=服务端清理残存房间

```json
{
    "seq": "创建房间的广播",
    "cmd": "1001",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1001",
            "timestamp": 1741331771,
            "room_Id": 0,
            "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323",
            "user_id": "20210",
            "turn": 0,
            "room_name": "'s lobby",
            "status": 1,
            "user_num_limit": 2,
            "room_type": 1,
            "room_level": 0,
            "room_user_list": [
                {
                    "user_id": "20210",
                    "nickname": "",
                    "turn": 1,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": true,
                    "is_ready": 0,
                    "is_my_turn": false,
                    "seat": 1,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                }
            ]
        }
    }
}
```


#  邀请好友 

```json
{
    "seq": "邀请好友 （user_id 是被邀请人）",
    "cmd": "mebInviteFriend",
    "service_id": "meme_battle",
    "data": {
        "user_id": "20211",
        "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323"
    }
}
```

```json
{
    "seq": "邀请好友",
    "cmd": "mebInviteFriend",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "邀请好友的广播（被邀请人会收到）",
    "cmd": "1015",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1015",
            "timestamp": 0,
            "user_id": "20211",
            "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323"
        }
    }
}
```

# === 加入房间
```json
{
    "seq": "加入房间",
    "cmd": "mebJoinRoom",
    "service_id": "meme_battle",
    "data": {
        "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323"
    }
}
```

```json
{
    "seq": "加入房间",
    "cmd": "mebJoinRoom",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "加入房间的广播",
    "cmd": "1003",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1003",
            "timestamp": 1741332253,
            "room_Id": 16,
            "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323",
            "user_id": "20211",
            "turn": 0,
            "room_name": "'s lobby",
            "status": 2,
            "user_num_limit": 2,
            "room_type": 1,
            "room_level": 0,
            "room_user_list": [
                {
                    "user_id": "20210",
                    "nickname": "",
                    "turn": 0,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": true,
                    "is_ready": 0,
                    "is_my_turn": false,
                    "seat": 1,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                },
                {
                    "user_id": "20211",
                    "nickname": "",
                    "turn": 0,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": false,
                    "is_ready": 0,
                    "priority_act": false,
                    "is_my_turn": false,
                    "seat": 2,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                }
            ]
        }
    }
}
```

# === 离开房间
```json
{
    "seq": "离开房间",
    "cmd": "mebLeaveRoom",
    "data": {
        "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323"
    }
}
```

```json
{
    "seq": "离开房间",
    "cmd": "mebLeaveRoom",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "离开房间的广播 （离开房间的人直接退出房间即可 ，只有在房间的人才能收到）",
    "cmd": "1005",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1005",
            "timestamp": 0,
            "user_id": "20211",
            "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323",
            "is_owner_leave": false,
            "new_owner": ""
        }
    }
}
```

# === 踢人

```json
    {
    "seq": "踢人 （被踢的人 user_id）",
    "cmd": "mebKickRoom",
    "data": {
      "user_id": "20211",
   	   "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323"
    }
  }
```

```json
{
    "seq": "踢人",
    "cmd": "mebKickRoom",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "踢人广播",
    "cmd": "1013",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1013",
            "timestamp": 0,
            "user_id": "20211",
            "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323"
        }
    }
}
```

# 获取用户当前游戏状态

```json
   {
    "seq": "获取用户当前游戏状态",
    "cmd": "mebUserState",
    "service_id": "meme_battle",
    "data": {
    }
  }
```

```json
{
    "seq": "获取用户当前游戏状态",
    "cmd": "mebUserState",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
           
        }
    }
}
```

```json
{
    "seq": "获取用户当前游戏状态的广播",
    "cmd": "1007",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1007",
            "timestamp": 1741333634,
            "user_id": "20210",
            "is_continue": true,
            "room_no": ""
        }
    }
}
```


# === 重新加入房间
```json
{
    "seq": "重新加入房间",
    "cmd": "mebReJoinRoom",
    "service_id": "meme_battle",
    "data": {
        "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323"
    }
}
```

```json
{
    "seq": "加入房间",
    "cmd": "mebReJoinRoom",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "加入房间的广播",
    "cmd": "1009",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1009",
            "timestamp": 1741332253,
            "room_Id": 16,
            "room_no": "51de2308-20dd-46e9-8225-3f2533fb5323",
            "user_id": "20211",
            "turn": 0,
            "room_name": "'s lobby",
            "status": 2,
            "user_num_limit": 2,
            "room_type": 1,
            "room_level": 0,
            "room_user_list": [
                {
                    "user_id": "20210",
                    "nickname": "",
                    "turn": 0,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": true,
                    "is_ready": 0,
                    "is_my_turn": false,
                    "seat": 1,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                },
                {
                    "user_id": "20211",
                    "nickname": "",
                    "turn": 0,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": false,
                    "is_ready": 0,
                    "priority_act": false,
                    "is_my_turn": false,
                    "seat": 2,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                }
            ]
        }
    }
}
```


# 开始游戏（房主）

```json
  {
    "seq": "开始对局游戏",
    "cmd": "mebStartPlay",
    "service_id": "meme_battle",
    "data": {
   	   "room_no": "127ea1da-cefc-4b90-8175-3b61179a13d6"
    }
  }

```

```json
{
    "seq": "开始对局游戏",
    "cmd": "mebStartPlay",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "开始游戏广播",
    "cmd": "1017",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1017",
            "timestamp": 1741612924,
            "room_no": "127ea1da-cefc-4b90-8175-3b61179a13d6",
            "room_user_list": [
                {
                    "user_id": "20210",
                    "nickname": "",
                    "turn": 1,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": true,
                    "is_ready": 1,
                    "seat": 1,
                    "is_my_turn": false,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                },
                {
                    "user_id": "20211",
                    "nickname": "",
                    "turn": 1,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": false,
                    "is_ready": 0,
                    "seat": 2,
                    "is_my_turn": false,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                }
            ]
        }
    }
}
```



# 加载完成（每个用户收到开始游戏的广播后 需要调用加载；最后一个加载完成 会发送问题描述）

```json

  {
    "seq": "加载",
    "cmd": "mebLoadCompleted",
    "service_id": "meme_battle",
    "data": {
   	   "room_no": "127ea1da-cefc-4b90-8175-3b61179a13d6"
    }
  }

```

```json
{
    "seq": "加载",
    "cmd": "mebLoadCompleted",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

[//]: # (game_status:游戏阶段)
[//]: # (	//0=游戏未开始)
[//]: # (	//1=游戏开始但是没有加载完成)
[//]: # (	//2=用户随牌阶段)
[//]: # (	//3=用户出牌阶段)
[//]: # (	//4=用户点赞阶段)
[//]: # (	//5=点赞界面 等待结算或者进入下一轮)

[//]: # (time_down:倒计时,这个字段需要配合game_status 使用，如下面game_status==3就是出牌阶段 ，time_down如果大于0 就是某个阶段倒计时，如果等于0，就是某个阶段已经结束，这个情况说明下个动作还没开始，直接等消息就行) 
[//]: # (就比如上一个例子 ，说明出牌倒计时已经结束，游戏状态要么处于服务端托管状态，但是托管（兜底）逻辑还没执行) 

```json
{
  "seq": "加载返回：出牌阶段",
  "cmd": "1019",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "proto_numb": "1019",
      "timestamp": 1742824104,
      "room_Id": 0,
      "room_no": "eb519b19-06d1-4c88-9065-9a00ebe1131b",
      "user_id": "111",
      "turn": 1,
      "room_name": "11119's lobby",
      "status": 4,
      "user_num_limit": 2,
      "room_type": 1,
      "room_level": 0,
      "time_down": 0,
      "game_status": 3,
      "curr_issue": {
        "issue_id": 6,
        "level": 0,
        "class": 0,
        "desc": "你正在走神，突然被老师点名回答问题！"
      },
      "room_user_list": [
        {
          "user_id": "111",
          "nickname": "",
          "turn": 1,
          "is_robot": 0,
          "is_leave": 0,
          "is_owner": true,
          "is_ready": 1,
          "seat": 1,
          "user_limit_num": 0,
          "user_cards": {
            "prev_user": "",
            "user_id": "111",
            "out_card_num": 0,
            "card_num": 4,
            "be_doubt_card_num": 0,
            "card_list": [
              {
                "card_id": 33,
                "card_type": 0,
                "card_level": 1,
                "name": "V1CA003",
                "point": 0,
                "express": 0,
                "suffix": "png",
                "img_url": "",
                "add_rate": 1
              },
              {
                "card_id": 33,
                "card_type": 0,
                "card_level": 1,
                "name": "V1CA003",
                "point": 0,
                "express": 0,
                "suffix": "png",
                "img_url": "",
                "add_rate": 1
              }
            ]
          },
          "win_price": 0,
          "bet": 0
        },
        {
          "user_id": "777",
          "nickname": "",
          "turn": 1,
          "is_robot": 0,
          "is_leave": 0,
          "is_owner": false,
          "is_ready": 0,
          "seat": 2,
          "user_limit_num": 0,
          "user_cards": {
            "prev_user": "",
            "user_id": "",
            "out_card_num": 0,
            "card_num": 0,
            "be_doubt_card_num": 0
          },
          "win_price": 0,
          "bet": 0
        }
      ],
      "other_user_cards": [
        {
          "prev_user": "",
          "user_id": "777",
          "out_card_num": 0,
          "card_num": 4,
          "be_doubt_card_num": 0
        }
      ]
    }
  }
}

```

```json
{
  "seq": "加载返回：点赞阶段",
  "cmd": "1019",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "proto_numb": "1019",
      "timestamp": 1742826726,
      "room_Id": 0,
      "room_no": "cc5069f2-5cb5-4465-9013-5ddf667897b5",
      "user_id": "111",
      "turn": 1,
      "room_name": "11119's lobby",
      "status": 4,
      "user_num_limit": 2,
      "room_type": 1,
      "room_level": 0,
      "curr_issue": {
        "issue_id": 42,
        "level": 0,
        "class": 0,
        "desc": "网购‘静音键盘’，敲击声比雷还响"
      },
      "time_down": 0,
      "game_status": 4,
      "room_user_list": [
        {
          "user_id": "111",
          "nickname": "",
          "turn": 1,
          "is_robot": 0,
          "is_leave": 0,
          "is_owner": true,
          "is_ready": 1,
          "seat": 1,
          "user_limit_num": 0,
          "user_cards": {
            "prev_user": "",
            "user_id": "111",
            "out_card_num": 1,
            "card_num": 3,
            "be_doubt_card_num": 0,
            "card_list": [
              {
                "card_id": 33,
                "card_type": 0,
                "card_level": 1,
                "name": "V1CA003",
                "point": 0,
                "express": 0,
                "suffix": "png",
                "img_url": "",
                "add_rate": 1
              },
              {
                "card_id": 18,
                "card_type": 0,
                "card_level": 1,
                "name": "V0CA018",
                "point": 0,
                "express": 0,
                "suffix": "png",
                "img_url": "",
                "add_rate": 1
              },
              {
                "card_id": 39,
                "card_type": 0,
                "card_level": 2,
                "name": "V1CA009",
                "point": 0,
                "express": 0,
                "suffix": "png",
                "img_url": "",
                "add_rate": 1.25
              }
            ]
          },
          "win_price": 0,
          "bet": 0
        },
        {
          "user_id": "777",
          "nickname": "",
          "turn": 1,
          "is_robot": 0,
          "is_leave": 0,
          "is_owner": false,
          "is_ready": 0,
          "seat": 2,
          "user_limit_num": 0,
          "user_cards": {
            "prev_user": "",
            "user_id": "",
            "out_card_num": 0,
            "card_num": 0,
            "be_doubt_card_num": 0
          },
          "win_price": 0,
          "bet": 0
        }
      ],
      "other_user_cards": [
        {
          "prev_user": "",
          "user_id": "777",
          "out_card_num": 1,
          "card_num": 3,
          "be_doubt_card_num": 0
        }
      ],
      "like_cards": [
        {
          "card_id": 17,
          "like_user_id": "777",
          "user_id": "111",
          "card_level": 1,
          "add_rate": 1,
          "like_num": 1
        }
      ],
      "out_carts": [
        {
          "card_id": 1,
          "card_type": 0,
          "card_level": 1,
          "name": "V0CA001",
          "point": 0,
          "express": 0,
          "suffix": "png",
          "img_url": "",
          "user_id": "111",
          "add_rate": 1
        },
        {
          "card_id": 14,
          "card_type": 0,
          "card_level": 1,
          "name": "V0CA014",
          "point": 0,
          "express": 0,
          "suffix": "png",
          "img_url": "",
          "user_id": "777",
          "add_rate": 1
        }
      ]
    }
  }
}
```


```json
{
    "seq": "问题描述（最后一个用户加载完成会发送）",
    "cmd": "1021",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1021",
            "timestamp": 1741613301,
            "issue": {
                "issue_id": 3,
                "level": 0,
                "class": 0,
                "desc": "你抢到了红包，但金额是最小的！"
            }
        }
    }
}
```

[//]: # (card_level //等级 1=流辉级 2=幻彩级 3=璀璨)

```json
{
  "seq": "用户手牌",
  "cmd": "1022",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "proto_numb": "1022",
      "timestamp": 1742826449,
      "user_id": "777",
      "room_no": "cc5069f2-5cb5-4465-9013-5ddf667897b5",
      "turn": 1,
      "cards": [
        {
          "card_id": 17,
          "card_type": 0,
          "card_level": 1,
          "name": "V0CA017",
          "point": 0,
          "express": 0,
          "suffix": "png",
          "img_url": "",
          "add_rate": 1
        },
        {
          "card_id": 10,
          "card_type": 0,
          "card_level": 2,
          "name": "V0CA010",
          "point": 0,
          "express": 0,
          "suffix": "png",
          "img_url": "",
          "add_rate": 1.25
        },
        {
          "card_id": 2,
          "card_type": 0,
          "card_level": 1,
          "name": "V0CA002",
          "point": 0,
          "express": 0,
          "suffix": "png",
          "img_url": "",
          "add_rate": 1
        },
        {
          "card_id": 4,
          "card_type": 0,
          "card_level": 1,
          "name": "V0CA004",
          "point": 0,
          "express": 0,
          "suffix": "png",
          "img_url": "",
          "add_rate": 1
        }
      ]
    }
  }
}
```


#  操作牌 （0:看牌 1:出牌 2:表情 3:重随）

[//]: # (type OperateCardReq struct {)
[//]: # (UserID  string  `json:"user_id"`)
[//]: # (EmojiId string  `json:"emoji_id"`)
[//]: # (RoomNo  string  `json:"room_no"`  //房间编号)
[//]: # (OpeType int     `json:"ope_type"` //0:看牌 1:出牌 2:表情 3:重随)
[//]: # (Pitch   float32 `json:"pitch"`)
[//]: # (Yaw     float32 `json:"yaw"`)
[//]: # (Looking bool    `json:"looking"` //看牌传)
[//]: # (Card    []*Card `json:"card_list"`)
[//]: # (})
```json
{
  "seq": "操作牌：随牌",
  "cmd": "mebOperateCard",
  "service_id": "meme_battle",
  "data": {
    "room_no": "7c9c4009-3ce2-4115-9bc5-0de14c4ea353",
    "ope_type": 3
  }
}
```

```json
{
  "seq": "操作牌：表情",
  "cmd": "mebOperateCard",
  "service_id": "meme_battle",
  "data": {
    "room_no": "7c9c4009-3ce2-4115-9bc5-0de14c4ea353",
    "emoji_id": "xxx"
  }
}
```

```json
{
  "seq": "操作牌：出牌",
  "cmd": "mebOperateCard",
  "service_id": "meme_battle",
  "data": {
    "room_no": "7c9c4009-3ce2-4115-9bc5-0de14c4ea353",
    "ope_type": 1,
    "card_list": [
      {
        "card_id": 5,
        "card_type": 2,
        "point": 13
      }
    ]
  }
}
```

```json
{
    "seq": "随牌成功",
    "cmd": "1027",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1027",
            "timestamp": 0,
            "user_id": "111",
            "room_no": "7c9c4009-3ce2-4115-9bc5-0de14c4ea353",
            "turn": 1,
            "cards": [
                {
                    "card_id": 13,
                    "card_type": 0,
                    "name": "V0CA013",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": ""
                },
                {
                    "card_id": 9,
                    "card_type": 0,
                    "name": "V0CA009",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": ""
                },
                {
                    "card_id": 7,
                    "card_type": 0,
                    "name": "V0CA007",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": ""
                },
                {
                    "card_id": 1,
                    "card_type": 0,
                    "name": "V0CA001",
                    "point": 0,
                    "express": 0,
                    "suffix": "png",
                    "img_url": ""
                }
            ]
        }
    }
}
```

```json
{
    "seq": "出牌成功",
    "cmd": "1026",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1026",
            "timestamp": 1741787192,
            "user_id": "111",
            "out_card_num": 1,
            "card_num": 3,
            "pitch": 0,
            "yaw": 0,
            "looking": false,
            "emoji_id": "",
            "card": [
                {
                    "card_id": 6,
                    "card_type": 0,
                    "name": "V0CA006",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": ""
                },
                {
                    "card_id": 11,
                    "card_type": 0,
                    "name": "V0CA011",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": ""
                },
                {
                    "card_id": 12,
                    "card_type": 0,
                    "name": "V0CA012",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": ""
                }
            ]
        }
    }
}
```

```json
{
    "seq": "收到此消息进入点赞页面",
    "cmd": "1028",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1028",
            "timestamp": 1741788550,
            "room_no": "85c28ec1-2e21-46bc-b421-2a515b07773b",
            "out_cards": [
                {
                    "card_id": 9,
                    "card_type": 0,
                    "name": "V0CA009",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": "",
                    "user_id": ""
                },
                {
                    "card_id": 3,
                    "card_type": 0,
                    "name": "V0CA003",
                    "point": 0,
                    "express": 0,
                    "suffix": "png",
                    "img_url": "",
                    "user_id": ""
                }
            ]
        }
    }
}
```

# 点赞

[//]: # (type CardLikeReq struct {)
[//]: # (LikeUserID string  `json:"like_user_id"` //被点赞的用户ID)
[//]: # (UserID     string  `json:"user_id"`      //用户ID)
[//]: # (RoomNo     string  `json:"room_no"`      //房间编号)
[//]: # (Card       []*Card `json:"card_list"`    //被点赞的牌)
[//]: # (})
```json
{
    "seq": "点赞",
    "cmd": "mebCardLike",
    "service_id": "meme_battle",
    "data": {
        "room_no": "0d40d016-943e-4eb6-95fd-d99b4fe1c103",
        "like_user_id": "222",
        "card_list": [
            {
                "card_id": 12
            }
        ]
    }
}
```

```json
{
    "seq": "点赞",
    "cmd": "mebCardLike",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "点赞广播",
    "cmd": "1029",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1029",
            "timestamp": 1741870752,
            "like_user_id": "222",
            "card": [
                {
                    "card_id": 12,
                    "card_type": 0,
                    "name": "V0CA012",
                    "point": 0,
                    "express": 0,
                    "suffix": "jpg",
                    "img_url": "",
                    "user_id": "222"
                }
            ]
        }
    }
}
```

```json
{
    "seq": "本局结束 计算本局最终结果的广播",
    "cmd": "1030",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1030",
            "timestamp": 1741870788,
            "room_no": "0d40d016-943e-4eb6-95fd-d99b4fe1c103",
            "like_detail_list": [
                {
                    "user_id": "222",
                    "nickname": "",
                    "head_photo": "222",
                    "on_go_like_num": 2,
                    "integral": 100,
                    "experience": 150
                },
                {
                    "user_id": "111",
                    "nickname": "",
                    "head_photo": "111",
                    "on_go_like_num": 2,
                    "integral": 100,
                    "experience": 0
                }
            ]
        }
    }
}
```


```json
{
    "seq": "游戏结束广播",
    "cmd": "1031",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1031",
            "timestamp": 1741870790,
            "room_Id": 0,
            "room_no": "0d40d016-943e-4eb6-95fd-d99b4fe1c103",
            "user_id": "",
            "turn": 0,
            "room_name": "",
            "status": 0,
            "user_num_limit": 0,
            "room_type": 0,
            "room_level": 0,
            "next_room_no": "bd4933ff-49b2-466a-bf5c-a2073ffcbe63",
            "room_user_list": [
                {
                    "user_id": "111",
                    "nickname": "",
                    "turn": 2,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": true,
                    "is_ready": 1,
                    "seat": 1,
                    "is_my_turn": false,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                },
                {
                    "user_id": "222",
                    "nickname": "",
                    "turn": 2,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": false,
                    "is_ready": 0,
                    "seat": 2,
                    "is_my_turn": false,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                }
            ]
        }
    }
}
```


# 加入匹配

```json
{
  "seq": "匹配开始",
  "cmd": "mebQuickMatchRoom",
  "service_id": "meme_battle",
  "data": {
    "room_no": "fc17c168-d841-4e91-b888-03161d14e50d"
  }
}
```

```json
{
    "seq": "匹配开始",
    "cmd": "mebQuickMatchRoom",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

```json
{
    "seq": "加入匹配队列成功的通知",
    "cmd": "1032",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1032",
            "timestamp": 1742267360,
            "room_Id": 0,
            "room_no": "",
            "user_id": "",
            "turn": 0,
            "room_name": "",
            "status": 0,
            "user_num_limit": 0,
            "room_type": 0,
            "room_level": 0
        }
    }
}
```

[//]: # (匹配逻辑流程是：单排为例子)
[//]: # (    1:每个用户分别调用创建房间接口)
[//]: # (    2：调用mebQuickMatchRoom 加入房间)
[//]: # (    3:等待匹配成功 会收到1034广播，证明匹配成功，（注意：就是匹配成功后的房间编号，和创建房间号不一样，匹配成功后，四个人进入一个新房间）)
[//]: # (    4:请求加载接口（匹配成功后直接请求加载，加入房间就绪等服务端已经处理）)

```json
{
    "seq": "匹配成功并开始；当收到这个消息，证明匹配成功，直接请求加载（注意就是匹配成功后，加入的房间和创建房间号不一样）",
    "cmd": "1034",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1034",
            "timestamp": 1742267559,
            "room_Id": 193,
            "room_no": "b121a51c-46b3-4cdd-a917-300892b734ef",
            "user_id": "111",
            "turn": 1,
            "room_name": "'s lobby",
            "status": 4,
            "user_num_limit": 4,
            "room_type": 2,
            "room_level": 0,
            "room_user_list": [
                {
                    "user_id": "111",
                    "nickname": "",
                    "turn": 1,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": true,
                    "is_ready": 1,
                    "seat": 1,
                    "is_my_turn": false,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                },
                {
                    "user_id": "222",
                    "nickname": "",
                    "turn": 1,
                    "is_robot": 0,
                    "is_killed": 0,
                    "is_leave": 0,
                    "is_owner": true,
                    "is_ready": 1,
                    "seat": 2,
                    "is_my_turn": false,
                    "user_limit_num": 0,
                    "win_price": 0,
                    "bet": 0,
                    "user_cards": {
                        "prev_user": "",
                        "user_id": "",
                        "out_card_num": 0,
                        "card_num": 0,
                        "be_doubt_card_num": 0
                    }
                }
            ]
        }
    }
}
```


# 取消匹配

```json
{
  "seq": "取消匹配",
  "cmd": "mebCancelMatchRoom",
  "service_id": "meme_battle",
  "data": {
    "room_no": "bfd616d9-acfe-4675-baaa-378e7d717fdb"
  }
}
```

```json
{
  "seq": "取消匹配",
  "cmd": "mebCancelMatchRoom",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {}
  }
}
```

```json
{
    "seq": "取消匹配广播",
    "cmd": "1033",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1033",
            "timestamp": 1742269338,
            "room_Id": 0,
            "room_no": "",
            "user_id": "",
            "turn": 0,
            "room_name": "",
            "status": 0,
            "user_num_limit": 0,
            "room_type": 0,
            "room_level": 0
        }
    }
}
```




# 就绪
```json
    {
    "seq": "准备就绪（就绪）",
    "cmd": "mebReady",
    "service_id": "meme_battle",
    "data": {
     "room_no":"1725f6e8-059c-4681-9a78-c2b0c36bc500"
    }
  }

```
```json
{
    "seq": "准备就绪（就绪）",
    "cmd": "mebReady",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```
```json
{
    "seq": "",
    "cmd": "1035",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1035",
            "timestamp": 0,
            "user_id": "222",
            "room_no": "1725f6e8-059c-4681-9a78-c2b0c36bc500"
        }
    }
}
```


# 取消就绪
```json
    {
    "seq": "取消就绪",
    "cmd": "mebCancelReady",
    "service_id": "meme_battle",
    "data": {
     "room_no":"1725f6e8-059c-4681-9a78-c2b0c36bc500"
    }
  }

```
```json
{
    "seq": "取消就绪",
    "cmd": "mebCancelReady",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```
```json
{
    "seq": "",
    "cmd": "1055",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "proto_numb": "1035",
            "timestamp": 0,
            "user_id": "222",
            "room_no": "1725f6e8-059c-4681-9a78-c2b0c36bc500"
        }
    }
}
```


# === 下面是接口  ===
[//]: # (调用逻辑同上 )
[//]: # (不同点是 没有广播，以返回结果为准)


# 版本列表
```json
{
  "seq": "版本列表",
  "cmd": "mebCardVersionList",
  "service_id": "meme_battle",
  "data": {
   
  }
}
```

```json
{
    "seq": "版本列表",
    "cmd": "mebCardVersionList",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "card_version_list": [
                {
                    "version": 1
                },
                {
                    "version": 2
                }
            ]
        }
    }
}
```

# 拆包 
```json
{
  "seq": "拆包 (num 拆包次数 测试环境如果拆10次，有时候可能会执行超时)",
  "cmd": "mebUnpackCard",
  "service_id": "meme_battle",
  "data": {
    "version":1,
    "num":1
  }
}
```

```json

{
  "seq": "拆包 (level 等级 1=流辉级 2=幻彩级 3=璀璨)",
  "cmd": "mebUnpackCard",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "list_card": [
        {
          "cards": [
            {
              "card_id": 32,
              "name": "V1CA002",
              "suffix": "png",
              "level": 1
            },
            {
              "card_id": 33,
              "name": "V1CA003",
              "suffix": "png",
              "level": 1
            },
            {
              "card_id": 37,
              "name": "V1CA007",
              "suffix": "png"
            },
            {
              "card_id": 40,
              "name": "V1CA010",
              "suffix": "png"
            },
            {
              "card_id": 39,
              "name": "V1CA009",
              "suffix": "png"
            }
          ]
        },
        {
          "cards": [
            {
              "card_id": 38,
              "name": "V1CA008",
              "suffix": "png"
            },
            {
              "card_id": 34,
              "name": "V1CA004",
              "suffix": "png"
            },
            {
              "card_id": 32,
              "name": "V1CA002",
              "suffix": "png"
            },
            {
              "card_id": 31,
              "name": "V1CA001",
              "suffix": "png"
            },
            {
              "card_id": 33,
              "name": "V1CA003",
              "suffix": "png"
            }
          ]
        }
      ]
    }
  }
}


```


# 图鉴列表 返回值字段 is_have_next_page == ture 代表有下一页，如果返回结构没有这个字段 说明is_have_next_page == false
```json
{
    "seq": "图鉴列表 (last_id 列表最后一个牌ID，分页用 ;level 等级 1=流辉级 2=幻彩级 3=璀璨)",
    "cmd": "mebHandbookList",
    "service_id": "meme_battle",
    "data": {
        "last_id": 0,
        "level": 0
    }
}
```

```json
{
  "seq": "图鉴列表 (last_id 列表最后一个牌ID，分页用 )",
  "cmd": "mebHandbookList",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "is_have_next_page": true,
      "all_cart_count": 28,
      "hand_list_card": [
        {
          "card_id": 31,
          "name": "V1CA001",
          "suffix": "png",
          "level": 1
        },
        {
          "card_id": 32,
          "name": "V1CA002",
          "suffix": "png",
          "is_own": true,
          "level": 1
        },
        {
          "card_id": 33,
          "name": "V1CA003",
          "suffix": "png",
          "is_own": true,
          "level": 1
        },
        {
          "card_id": 34,
          "name": "V1CA004",
          "suffix": "png"
        },
        {
          "card_id": 35,
          "name": "V1CA005",
          "suffix": "png",
          "is_own": true
        }
      ]
    }
  }
}
```


# 好友列表 返回值字段 is_have_next_page == ture 代表有下一页，如果返回结构没有这个字段 说明is_have_next_page == false
```json
  {
    "seq": "好友列表",
    "cmd": "mebFriendList",
    "service_id": "meme_battle",
    "data": {
       "last_id": 0
    }
  }

```
```json
{
  "seq": "好友列表",
  "cmd": "mebFriendList",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "is_have_next_page": true,
      "user_friend": [
        {
          "friend_user_id": "222",
          "nickname": "222_敬请期待"
        }
      ]
    }
  }
}
```

# 申请好友列表 返回值字段 is_have_next_page == ture 代表有下一页，如果返回结构没有这个字段 说明is_have_next_page == false
```json
  {
    "seq": "申请好友列表",
    "cmd": "mebAuditUserList",
    "service_id": "meme_battle",
    "data": {
       "last_id": 0
    }
  }

```
```json
{
  "seq": "申请好友列表",
  "cmd": "mebAuditUserList",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "audit_user": [
        {
          "application_user": "333",
          "nickname": "333_敬请期待",
          "audit_id": 1
        },
        {
          "application_user": "444",
          "nickname": "444_敬请期待",
          "audit_id": 2
        }
      ]
    }
  }
}
```


#  添加好友
```json
  {
    "seq": "添加好友 输入的用户ID是审核人",
    "cmd": "mebAddFriend",
    "service_id": "meme_battle",
    "data": {
       "audit_user": "111"
    }
  }

```
```json
{
    "seq": "添加好友",
    "cmd": "mebAddFriend",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

#  审核好友
```json
  {
  "seq": "审核好友",
  "cmd": "mebAuthFriend",
  "service_id": "meme_battle",
  "data": {
    "audit_id" : 0
  }
}
```
```json
{
  "seq": "审核好友",
  "cmd": "mebAuthFriend",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {}
  }
}
```

#  删除好友
```json
    {
    "seq": "删除好友",
    "cmd": "mebDelFriend",
    "service_id": "meme_battle",
    "data": {
     "friend_id": 5
    }
  }

```
```json
{
    "seq": "删除好友",
    "cmd": "mebDelFriend",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {}
    }
}
```

#  用户资料
```json
  {
    "seq": "资料",
    "cmd": "mebUserDetail",
    "service_id": "meme_battle",
    "data": {
      
    }
  }
```
```json
{
  "seq": "资料",
  "cmd": "mebUserDetail",
  "response": {
    "code": 200,
    "codeMsg": "Success",
    "data": {
      "user_id": "111",
      "nickname": "11119"
    }
  }
}
```
```json
```


#  用户的经验值和金币值
```json
  {
    "seq": "用户的经验值和金币值",
    "cmd": "mebGetCoinExperience",
    "service_id": "meme_battle",
    "data": {
      
    }
  }
```
```json
{
    "seq": "用户的经验值和金币值",
    "cmd": "mebGetCoinExperience",
    "response": {
        "code": 200,
        "codeMsg": "Success",
        "data": {
            "coin_num": 11,
            "experience": 1
        }
    }
}
```


```json
```
```json
```