FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.4

COPY . .
RUN swag init --output internal/http/swagger/docs --generalInfo cmd/main.go
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH=amd64

COPY . .

RUN go build -tags musl -o mms-app -ldflags="-s -w -extldflags -static" ./cmd

FROM golang:1.19-alpine

WORKDIR /

COPY --from=build /app/mms-app /

EXPOSE 8080

ENTRYPOINT [ "./mms-app" ]
