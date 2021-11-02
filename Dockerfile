FROM alpine:latest

LABEL maintainer="1272105563@qq.com"

RUN apk add --update --no-cache tzdata

RUN mkdir -p "/app"

COPY wiblog /app/wiblog
COPY conf /app/conf
COPY website /app/website
COPY assets /app/assets

EXPOSE 9000

WORKDIR /app
CMD ["./wiblog"]