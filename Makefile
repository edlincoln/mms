#performs the configuration of environment properties and installation of important dependencies such as swaggo
#which is responsible for generating the documentation, mockery which performs the mock up of services and repositories.
install:
	go env -w GOPROXY=https://proxy.golang.org
	go env -w CGO_ENABLED="1"
	go env -w GO111MODULE='on'
	go env -w GOBIN=${GOPATH}/bin
	go mod download -x
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@v1.8.4
	swag -v
	go install github.com/vektra/mockery/v2@latest

#run the application
run:
	go run cmd/main.go

#run the application
run-install: install doc-swag
	go run cmd/main.go

#download project dependencies
dep:
	go get -d -u ./...
	go mod download -x

doc-swag:
	swag init --parseDependency --parseInternal --parseDepth 1 --output internal/http/swagger/docs --generalInfo cmd/main.go

#build the application
build:
	go build -race -o bin/${BINARY_NAME} -ldflags="-s -w" ./cmd

#configura e executa o docker-compose
docker-compose-up:
	docker compose -f ./resources/docker/docker-compose.yml up -d

docker-compose-down:
	docker compose -f ./resources/docker/docker-compose.yml down