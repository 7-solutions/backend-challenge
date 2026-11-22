MONGODB_NAME ?= backend-challenge

.PHONY: mockgen-users-repo mockgen-users-usecase mockgen-auth-usecase mockgen-all swag seed-admin

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
