FROM golang:1.15.4-buster as builder

ARG APP_CMD_DIR

RUN apt-get update && apt-get install -y libc6 curl wget

COPY build/keys/id_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig && \
    echo "StrictHostKeyChecking no " > /root/.ssh/config
RUN go env -w GOPRIVATE=github.com/kettari/*

COPY . /src/github.com/kettari/shitdetector

RUN cd /src/github.com/kettari/shitdetector && make PROJECT_CMD=$APP_CMD_DIR

FROM debian:10

RUN apt-get update && apt-get install -y locales

# Set the locale
RUN locale-gen ru_RU.UTF-8
ENV LANG ru_RU.UTF-8
ENV LANGUAGE ru_RU:ru
ENV LC_ALL ru_RU.UTF-8

RUN rm /etc/localtime && ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime

COPY --from=builder /src/github.com/kettari/shitdetector/bin/shitdetector_bot /var/www/shitdetector_bot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /var/www

CMD ["/var/www/shitdetector_bot"]
