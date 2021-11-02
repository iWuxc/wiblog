FROM golang:1.17

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release \
    PORT=9000

LABEL maintainer="1272105563@qq.com"

#RUN apk add --update --no-cache tzdata

RUN mkdir -p "/app"

COPY wiblog /app/wiblog
COPY conf /app/conf
COPY website /app/website
COPY assets /app/assets

EXPOSE 9000

WORKDIR /app
CMD ["./wiblog"]