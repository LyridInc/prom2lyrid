FROM golang:latest
WORKDIR /go/src/prom2lyrid/
COPY . .
ENV GO111MODULE=on
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init
RUN go mod tidy
RUN CGO_ENABLED=0 go build -v -o app

FROM alpine
RUN apk add --no-cache ca-certificates bash
COPY --from=0 /go/src/prom2lyrid/app .
COPY --from=0 /go/src/prom2lyrid/docs ./docs
COPY --from=0 /go/src/prom2lyrid/.env .
ENTRYPOINT ["/app"]
