FROM quay.io/goswagger/swagger as base

WORKDIR /go/src/github.com/blinkops/blink-go-cli

COPY . .

RUN make gen
RUN make build

##################
FROM alpine:3 AS blink

WORKDIR /blink

COPY --from=base /go/src/github.com/blinkops/blink-go-cli/blink .

ENTRYPOINT ["./blink"]