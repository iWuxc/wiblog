FROM alpine:latest

LABEL maintainer="1272105563@qq.com"

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk add --update --no-cache tzdata

RUN mkdir -p "/app"

RUN echo $(pwd)

COPY bin/wiblog /app/wiblog
COPY conf /app/conf
COPY website /app/website
COPY assets /app/assets

EXPOSE 9000

WORKDIR /app
CMD ["./wiblog"]