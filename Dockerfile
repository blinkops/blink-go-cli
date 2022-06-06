FROM quay.io/goswagger/swagger as base

WORKDIR /go/src/github.com/blinkops/blink-go-cli

COPY . .

RUN make gen
RUN make build

##################
FROM alpine:3.15 AS blink-cli

WORKDIR /blink-cli

COPY --from=base /go/src/github.com/blinkops/blink-go-cli/blink-cli .

ENTRYPOINT ["./blink-cli"]