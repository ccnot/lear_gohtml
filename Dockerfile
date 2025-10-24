FROM alpine:3.14.2
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# 获取 需要的依赖项。
RUN apk add --no-cache openssl openssl-dev nghttp2-dev ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /godash
COPY ./main .
COPY ./config ./config
COPY ./web ./web

EXPOSE 80
ENTRYPOINT ["./main"]