FROM golang:1.13.3 as builder
ARG BIN
ENV BIN_CMD_DIR ${BIN}

COPY . /go/src/lionnix-metrics-api
WORKDIR /go/src/lionnix-metrics-api

RUN make native BIN=$BIN_CMD_DIR
RUN mkdir /app && cp -r dist/linux-amd64/$BIN_CMD_DIR/* /app

FROM alpine:3.10
RUN apk --no-cache add ca-certificates rsync openssh
WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

# expose port for api
EXPOSE 8888

ENTRYPOINT ["./run.sh", "start"]
