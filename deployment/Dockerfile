FROM golang:1.24-alpine3.21 as builder

RUN mkdir -p /go/src/github.com/abdussalamfaqih/wallet-service-dev


WORKDIR /go/src/github.com/abdussalamfaqih/wallet-service-dev

COPY . .

RUN apk add build-base

RUN go mod tidy -compat=1.21 && go mod download && go mod vendor
RUN CGO_ENABLED=1 GOOS=linux go build -a .

FROM alpine:3.21

COPY --from=builder /go/src/github.com/abdussalamfaqih/wallet-service-dev .

RUN chmod +x ./wallet-service-dev

USER root

# executable
ENTRYPOINT [ "./wallet-service-dev" ]

# # arguments that can be overridden
CMD [ "run-http" ]