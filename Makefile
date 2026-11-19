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