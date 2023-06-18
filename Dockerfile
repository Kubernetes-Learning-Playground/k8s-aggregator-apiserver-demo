FROM golang:1.18.7-alpine3.15 as builder

WORKDIR /app

# copy modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# cache modules
RUN go mod download

# copy source code
COPY pkg/ pkg/
COPY cmd/ cmd/
COPY cert/ cert/
COPY resources/ resources/
# build
RUN CGO_ENABLED=0 go build \
    -a -o kube-aggregator-apiserver cmd/main.go

FROM alpine:3.13
WORKDIR /app

USER nobody

COPY --from=builder --chown=nobody:nobody /app/kube-aggregator-apiserver .
# 加载需要的证书文件
COPY --from=builder /app/cert cert/
# FIXME 暂时使用项目中复制下来的kube-config文件，从集群中复制下来的(正确的应该使用挂载)
COPY --from=builder /app/resources resources/

ENTRYPOINT ["./kube-aggregator-apiserver"]