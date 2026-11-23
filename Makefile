MONGODB_NAME ?= backend-challenge
PROTOC_GEN_GO_VERSION ?= v1.36.12
PROTOC_GEN_GO_GRPC_VERSION ?= v1.5.1
MOCKGEN_VERSION ?= v0.6.0
SWAG_VERSION ?= v1.16.6
PROTO_FILES := proto/users/v1/user_service.proto
GO_PACKAGES := $(shell go list ./... | grep -vE '/app/(gen|docs|test_helpers|mocks)(/|$$)')

.PHONY: dev test coverage install-tools install-protobuf-tools mockgen-users-repo mockgen-users-usecase mockgen-auth-usecase mockgen-all swag seed-admin gen-protobuf up log

dev: install-tools
	go run ./app/cmd/server

test:
	go test $(GO_PACKAGES)

coverage:
	go test -coverprofile=coverage.out $(GO_PACKAGES)
	go tool cover -func=coverage.out

install-tools: install-protobuf-tools
	go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)
	go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)


mockgen-users-repo:
	@echo "Generating mock for users-repo"
	mockgen -source=app/layers/repositories/users/repo.go -destination=app/mocks/users/repo.go -package=users

mockgen-users-usecase:
	@echo "Generating mock for users-usecase"
	mockgen -source=app/layers/usecases/users/usecase.go -destination=app/mocks/users/usecase.go -package=users

mockgen-auth-usecase:
	@echo "Generating mock for auth-usecase"
	mockgen -source=app/layers/usecases/auth/usecase.go -destination=app/mocks/auth/usecase.go -package=auth

mockgen-all: mockgen-users-repo mockgen-users-usecase mockgen-auth-usecase

swag:
	@echo "Generating swagger"
	swag init -g app/cmd/server/main.go --output app/docs --parseDependency --parseInternal

seed-admin:
	# secret123
	@echo "Seeding admin user into MongoDB"
	docker compose exec -T mongo mongosh --quiet $(MONGODB_NAME) --eval 'db.users.updateOne({email: "admin@admin.com"}, {$$setOnInsert: {name: "admin", email: "admin@admin.com", password: "$$2a$$10$$wdEndUaQxpH5Tk10Ig2uUueRMqJ1HTb4FdIIpt2af1TQfhnBU1uv2", created_at: new Date(), updated_at: new Date()}}, {upsert: true})'

install-protobuf-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

gen-protobuf: install-protobuf-tools
	PATH="$(shell go env GOPATH)/bin:$(PATH)" protoc --go_out=. --go_opt=module=github.com/chaihanij/backend-challenge \
		--go-grpc_out=. --go-grpc_opt=module=github.com/chaihanij/backend-challenge \
		$(PROTO_FILES)

up:
	docker compose up --build -d

log:
	docker compose logs app -f
