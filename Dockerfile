FROM ubuntu:22.04

USER root

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    curl \
    && curl -O https://dl.google.com/go/go1.21.6.linux-amd64.tar.gz \
    && tar xvf go1.21.6.linux-amd64.tar.gz \
    && mv go /usr/local \
    && rm -rf /var/lib/apt/lists 

ENV PATH=$PATH:/usr/local/go/bin
ENV TELEGRAM_APITOKEN=6764971311:AAGEp9hVEe2NkCCWxbnbu-s3_zQa8hRPaS0

COPY src /app

WORKDIR /app

RUN go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5 \
    && go get -u gorm.io/gorm \
    && go get -u gorm.io/driver/postgres \
    && go get -u github.com/redis/go-redis/v9 

CMD ["go", "run", "."]
