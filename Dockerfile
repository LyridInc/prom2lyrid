FROM golang:alpine
RUN apk add --no-cache nodejs npm
WORKDIR /go/src/prom2lyrid/
COPY . .
ENV GO111MODULE=on
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init
RUN go mod tidy
RUN CGO_ENABLED=0 go build -v -o app

WORKDIR /go/src/prom2lyrid/web
RUN npm install
RUN npm run build

FROM alpine
RUN apk add --no-cache ca-certificates bash
WORKDIR /prom2lyrid/
COPY --from=0 /go/src/prom2lyrid/app .
COPY --from=0 /go/src/prom2lyrid/docs ./docs
COPY --from=0 /go/src/prom2lyrid/.env .
COPY --from=0 /go/src/prom2lyrid/web/build ./web/build
ENTRYPOINT ["/prom2lyrid/app"]
