FROM registry.cn-shenzhen.aliyuncs.com/xxx/golang:1.19.1 AS build-env

WORKDIR /go/src/codeup.aliyun.com/devops/consul/consul-webhook

ADD . /go/src/codeup.aliyun.com/devops/consul/consul-webhook

RUN go build -o main

FROM registry.cn-shenzhen.aliyuncs.com/xxx/alpine-hzjy:1.1.0

WORKDIR /go

COPY --from=build-env /go/src/codeup.aliyun.com/devops/consul/consul-webhook/main ./main
COPY --from=build-env /go/src/codeup.aliyun.com/devops/consul/consul-webhook/docker/ docker/

# 容器内nobody为65534，宿主机可能为centos 99  rocky os为 65534
RUN chown -R 65534:65534 /go

EXPOSE 8080

ENTRYPOINT ["docker/entrypoint.sh"]
CMD ["./main"]
