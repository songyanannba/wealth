#!/usr/bin/env bash

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    im_protobuf.proto



//===
#protoc --go_out=../ ./proto/game/*.proto
#protoc --go_out=../ --go_opt=paths=source_relative --go-grpc_out=../ --go-grpc_opt=paths=source_relative  ./proto/game/rat_mining.proto
#protoc --go_out=../ --go_opt=paths=source_relative --go-grpc_out=../ --go-grpc_opt=paths=source_relative  ./proto/rat/*.proto


/Users/syn/goProjects/zhuoshuo/nb_game_server/gate_way/pbs/rat/