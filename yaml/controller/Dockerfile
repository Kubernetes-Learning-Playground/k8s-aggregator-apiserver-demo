# 注意：需要移到项目根目录执行，本目录执行会报错。
# 此为controller项目的Dockerfile文件
# 但与api-server的Dockerfile冲突，因此放在此目录下
# Build the manager binary
FROM golang:1.18.7-alpine3.15 as builder

WORKDIR /app
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# 坑：
# 报错 go mod download: google.golang.org/api@v0.44.0: read tcp 172.17.0.3:60862->14.204.51.154:443: read: connection reset by peer
# The command '/bin/sh -c go mod download' returned a non-zero code: 1
# make: *** [docker-build] 错误 1
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
# # 需要把该放入的都copy进去，如果报出 package xxxxx is not in GOROOT  => 就是这个问题。
COPY cmd/controller/main.go cmd/controller/main.go
COPY pkg/ pkg/
# FIXME: 目前还是使用文件内挂载kubeconfig文件的方式，需要改成读内部sa
COPY resources/ resources/
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o myingresscontroller cmd/controller/main.go


FROM alpine:3.12
WORKDIR /app
# 需要的文件需要复制过来
COPY --from=builder /app/myingresscontroller .
COPY --from=builder /app/resources resources/
USER 65532:65532

ENTRYPOINT ["./myingresscontroller"]