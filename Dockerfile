FROM golang:1.22 as build-deps

WORKDIR /usr/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build /usr/src/app/cmd/main.go

FROM alpine:3.19.1
WORKDIR /usr/src/app

ARG MODE
COPY --from=build-deps /usr/src/app/run.sh run.sh
COPY --from=build-deps /usr/src/app/main main
COPY --from=build-deps /usr/src/app/configs/$MODE ./config

ARG MODULE
RUN ls /usr/src/app/config
RUN chmod +x run.sh
RUN apk add --no-cache bash
RUN apk add --no-cache libc6-compat

ENV LOG_PATH=/logs/$MODULE.log

ENTRYPOINT ["./run.sh"]
