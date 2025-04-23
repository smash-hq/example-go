# 运行阶段
FROM golang:alpine AS builder
# 设置工作目录
WORKDIR /build
# 复制源
COPY . .
ENTRYPOINT [ "sh", "-c", "cd /build && go run main.go" ]