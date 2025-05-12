package websocket

// OutCardCountDownTimeInt 生产环境 15（随牌） + 15（出牌） + 10（服务端滞后）
const OutCardCountDownTimeInt = 40

// CommTimeOut 通用超时时间 （重置牌倒计时 点赞倒计时 出牌倒计时）
const CommTimeOut = 15
const CommTimeOutDouble = 30

//const CommTimeOut = 60
//const CommTimeOutDouble = 120

const CommTimeDelay = 10

const RoomAlive = 60 * 20 //房间存活时间设置
