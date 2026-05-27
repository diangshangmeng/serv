#!/bin/bash

# 内存限制（软限制1GB，硬限制1.1GB）
export GOMEMLIMIT=1GiB
export GOGC=50                          # 降低GC阈值，更频繁回收内存
export GOMAXPROCS=1                     # 1核CPU设置为1

# go build -ldflags="-s -w" -o api-server cmd/main.go

# 启动服务
./api-server
