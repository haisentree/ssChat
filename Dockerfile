# ==================== 构建阶段 ====================
FROM golang:1.23-alpine AS builder

WORKDIR /build

# 先复制依赖清单，利用 Docker 层缓存：依赖不变时跳过 go mod download
# GOPROXY 使用国内代理，避免容器内访问 proxy.golang.org 超时
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

# 再复制源码并编译
# CGO_ENABLED=0 产出纯静态二进制，可在 alpine/scratch 中直接运行
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ssChat ./cmd

# ==================== 运行阶段 ====================
FROM alpine:3.20

WORKDIR /app

# 复制二进制和静态资源
# 程序用相对路径 ../internal/web 读取静态文件，
# 因此保持 与原始镜像相同的 /app/cmd + /app/internal 目录结构
COPY --from=builder /build/ssChat ./cmd/ssChat
COPY internal/web ./internal/web

# 非 root 用户运行
RUN adduser -D -u 10001 sschat && chown -R sschat:sschat /app
USER sschat

# 程序用相对路径 ../internal/web 读取静态文件，
# 工作目录必须是 /app/cmd，../internal/web 才能解析到 /app/internal/web
WORKDIR /app/cmd

EXPOSE 8080

ENTRYPOINT ["/app/cmd/ssChat"]
