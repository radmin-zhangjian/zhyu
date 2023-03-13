FROM golang:1.18-alpine3.16 AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build
#RUN adduser -u 10001 -D app-zhyu

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
ENV GOPROXY https://goproxy.cn
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .
COPY ./website /app/website
COPY ./runtime /app/runtime
COPY ./resources /app/resources
COPY ./setting/configFile /app/setting/configFile

# 将我们的代码编译成二进制可执行文件
RUN go build -o /app/app_zhyu .
#RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o /app/app_zhyu .

# 需要运行的命令 到这里就可以结束了 也可以用下面的scratch构建
#USER app-zhyu
#ENTRYPOINT ["./app_zhyu", "-config.file", "/configFile/app.yaml"]
#ENTRYPOINT ["/app/app_zhyu", "-config.file", "/configFile/app.yaml"]

###################
# 接下来创建一个小镜像
###################
FROM scratch

# 从builder镜像中把/dist/app 拷贝到当前目录
WORKDIR /app
COPY --from=builder /app/app_zhyu /app/
COPY --from=builder /app/setting/configFile /build/setting/configFile

# 需要运行的命令
#USER app-zhyu
ENTRYPOINT ["./app_zhyu", "-config.file", "configFile/app.yaml"]
