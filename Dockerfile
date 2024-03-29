FROM node:14-alpine3.14 as node-builder

COPY ./web /build
WORKDIR /build

RUN npm install \
    && npm run build

FROM golang:1.17-alpine3.14 as go-builder

COPY . /build
WORKDIR /build

COPY --from=node-builder /build/dist ./web/dist

ARG GOPROXY="https://proxy.golang.org"

ENV GO111MODULE=on
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/cloudrain

FROM alpine:3.14

WORKDIR /app

COPY --from=go-builder /build/cloudrain .

VOLUME /app/configs
VOLUME /app/tmp

ENTRYPOINT ["/app/cloudrain"]
CMD ["serve"]
