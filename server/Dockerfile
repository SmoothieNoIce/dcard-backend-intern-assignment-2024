FROM golang:1.20.1-bullseye AS build
RUN apt-get update -y
RUN apt-get install curl -y
RUN apt install graphviz -y
RUN apt install vim -y
RUN apt install lsof -y

# 安裝 gomigrate
RUN cd /tmp \
    && mkdir -p gomigrateDownload \
    && cd gomigrateDownload \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate $GOPATH/bin

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/cweill/gotests/gotests@latest
RUN go install github.com/fatih/gomodifytags@latest
RUN go install github.com/josharian/impl@latest
RUN go install github.com/haya14busa/goplay/cmd/goplay@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

WORKDIR /usr/src/dcard-backend-intern-assignment-2024/server