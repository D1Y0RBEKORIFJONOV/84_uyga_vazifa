
SWAGGER_CMD = swag
GO_RUN_CMD = go run

swagger-init:
	$(SWAGGER_CMD) init -g internal/http/handler/user.go

run:
	$(GO_RUN_CMD) cmd/app/main.go

all: swagger-init run
