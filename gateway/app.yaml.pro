app:
  logFile: log/gin.log
  httpPort: 8090
  webSocketPort: 8099
  rpcPort: 9001
  oreRpcUrl: 0.0.0.0:9001
  httpUrl: 0.0.0.0:8090
  webSocketUrl: 0.0.0.0:8099
  svName : gate_way
  environment : pro
mysql:
  path: rm-bp1sa3z6fmf3iiasr.rwlb.rds.aliyuncs.com
  port: "3306"
  config: "charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
  db-name: "go_game"
  username: "songyanan"
  password: "uV1%dM4#uQ6$xW3*nF3#mT7&uH0is1"
  max-idle-conns: 10
  max-open-conns: 100
  prefix: ""
  singular: false
  log-zap: false
redis:
  addr: r-bp1rvw2xzyk3gvsg8ipd.redis.rds.aliyuncs.com:6379
  password: 84*bQ8zGoS3@jJ5&aW9#lBjD!7sS@1
  DB: 0
  poolSize: 30
  minIdleConns: 30
zap:
  level: debug
  prefix: ''
  format: console
  director: log
  encode-level: LowercaseColorLevelEncoder
  stacktrace-key: stacktrace
  max-age: 30
  show-line: true
  log-in-console: true
huanCangUrl:
  getUserScoreUrl :  http://yishu-api.huancang.in/api/getaway_extra/score
  updateUserScoreUrl : http://yishu-api.huancang.in/api/getaway_extra/settle
  addUserScoreUrl :  http://yishu-api.huancang.in/api/bscore/operate
  getUserInfo: http://yishu-api.huancang.in/
jwt:
  signing-key: xxxx
  expires-time: 7d
  buffer-time: 1d
  issuer: qmPlus
db-list:
  - disable: false
    type: "mysql"
    alias-name: "read"
    path: ""
    port: "3306"
    #config: "charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
    config: "charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
    db-name: "go_game"
    username: "go_game"
    password: "n5ZTHJdWhTamEam2"
    prefix: ""
    singular: false
    engine: ""
    max-idle-conns: 10
    max-open-conns: 100
    log-mode: warn
    log-zap: false
cluster:
  gate-port: 34567
  etcd-endpoints: 127.0.0.1:2379
  etcd-prefix: ""
  etcd-user: "root"
  etcd-pass: "123456"
  nats-endpoints: 127.0.0.1:4222
  nats-user: "root"
  nats-pass: "123456"
system:
  env: public
  addr: 8888
  db-type: mysql
  oss-type: local
  use-multipoint: false
  use-redis: true
  api-log: false
  iplimit-count: 15000
  iplimit-time: 3600
  router-prefix: ""
  listen-ip: "0.0.0.0"
  connect-ip: "0.0.0.0"
  api-addr: "9888"
  master-addr: "34567"
  game-addr: "34580"
  bind-addr: "34570"
  gate-addr: "34590"
  backend-addr: "34500"
  migrate: true
  connect-cluster: false
  api-domain: ""
  game-domain: "xxx"
  storage-domain: ""
  ws-path: ""
  test-api-url: ""
  clusters:
    - name: "server"
      ip: "0.0.0.0"
      ws-scheme: "ws"
cors:
  mode: strict-whitelist
  whitelist:
    - allow-origin: example1.com
      allow-methods: POST, GET
      allow-headers: Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id
      expose-headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,
        Content-Type
      allow-credentials: true
    - allow-origin: example2.com
      allow-methods: GET, POST
      allow-headers: content-type
      expose-headers: Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,
        Content-Type
      allow-credentials: true
keys:
  ip-api-key: ""
timer:
  start: true
  spec: '@daily'
  with_seconds: false
  detail:
#    - tableName: sys_operation_records
#      compareField: created_at
#      interval: 720h
#    - tableName: jwt_blacklists
#      compareField: created_at
#      interval: 168h







