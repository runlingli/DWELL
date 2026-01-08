# # =========================
# # 第一阶段：构建 Go 二进制文件
# # =========================

# # 使用官方的 Go 镜像作为“构建环境”
# # alpine 版本体积更小，适合容器
# FROM golang:1.18-alpine as builder

# # 在容器中创建一个 /app 目录
# # 用来存放你的项目代码
# RUN mkdir /app
# # RUN：构建镜像时执行
# # CMD：容器启动时执行

# # 将当前宿主机目录下的所有文件
# # 拷贝到容器内的 /app 目录
# COPY . /app

# # 设置工作目录
# # 后续的命令都会在 /app 下执行
# WORKDIR /app

# # 编译 Go 程序
# # CGO_ENABLED=0 表示禁用 CGO，生成“纯静态二进制文件”
# # -o brokerApp 指定输出的可执行文件名
# # ./cmd/api 是你的 main 包所在路径
# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api


# # 给生成的二进制文件增加可执行权限
# # 确保在运行阶段可以被直接执行
# RUN chmod +x /app/brokerApp


# =========================
# 第二阶段：运行 Go 程序
# =========================

# 使用一个极小的 alpine 镜像作为运行环境
# 这里不再需要 Go 编译器
FROM alpine:latest

# 在运行镜像中创建 /app 目录
RUN mkdir /app

# 从第一阶段（builder）中
# 只拷贝已经编译好的二进制文件
# COPY --from=builder /app/brokerApp /app
COPY brokerApp /app

# 容器启动时执行的命令
# 直接运行 Go 编译好的二进制程序
CMD [ "/app/brokerApp" ]
