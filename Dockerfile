FROM node:18 AS app-builder

ADD plugin /app
WORKDIR /app

RUN npm i -g pnpm
RUN pnpm i && pnpm build

FROM golang:1.21-alpine AS server-builder

RUN apk add upx
RUN apk add -U --no-cache ca-certificates

WORKDIR /go/src/github.com/slntopp/core-chatting

COPY go.mod go.sum cc.yaml ./
ADD api api
ADD pkg pkg
ADD cc cc

RUN ls -la
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./api
RUN upx ./server

FROM scratch

COPY --from=app-builder /app/dist/ /dist
COPY --from=server-builder /go/src/github.com/slntopp/core-chatting/server /server
COPY --from=server-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

LABEL org.opencontainers.image.source https://github.com/slntopp/core-chatting

ENTRYPOINT ["/server"]