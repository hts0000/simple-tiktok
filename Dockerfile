# 启动编译环境
FROM golang:1.19

# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源码
COPY . /go/src/simple-tiktok

# 编译
WORKDIR /go/src/simple-tiktok
RUN go install simple-tiktok ./

# 暴露端口
EXPOSE 8080

# 设置服务入口
ENTRYPOINT [ "/go/bin/simple-tiktok" ]