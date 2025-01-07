ARG GO_VERSION=1.22.1

FROM golang:${GO_VERSION} AS build

WORKDIR /opt/zurl

COPY . ./

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zurl cmd/zurl/main.go

FROM istio/distroless:latest AS run

WORKDIR /opt/zurl

COPY --from=build /opt/zurl/zurl .

EXPOSE 8080

CMD ["./zurl"]