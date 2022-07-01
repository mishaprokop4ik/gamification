BIN_NAME := 'acheer'
GOLINT := golangci-lint
MIGRATE_BIN := migrate
POSTGRES_NAME ?= postgres

MIGRATE=migrate -path internal/persistence/postgres/migrations -database postgres://postgres:12345@localhost:5432/${POSTGRES_NAME}?sslmode=disable
version = 0.0.1-beta

swagger: # Generate swagger documentation
	 swag init --parseDependency --parseInternal -g ./cmd/acheer/main.go
dep: # Download required dependencies
	GOPRIVATE=${GOPRIVATE} go mod tidy
	GOPRIVATE=${GOPRIVATE} go mod download
	GOPRIVATE=${GOPRIVATE} go mod vendor

build: ## Build the binary file
	CGO_ENABLED=1 go build -o ./bin/${BIN_NAME} -a .

test: dep test-db-prepare ## Run unit tests
	go test -cover -v -race -count=1 -tags integration ./...

test-run-one: dep test-db-prepare ## Run unit tests
	go test -cover -v -race -run $(NAME) -count=1 -tags integration ./...

lint: dep check-lint ## Lint the files local env
	$(GOLINT) run --timeout=5m -c .golangci.yml

cilint: dep check-lint
	mkdir reports || true
	$(GOLINT) run --timeout=5m -c .golangci.yml ./... > reports/gometalinter-report.out
	cat reports/gometalinter-report.out

citest: dep ## Run unit tests ci env
	mkdir reports || true
	go test -race -count=1 -coverprofile="reports/test-report.out" ./...

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.0

check-migrate:
	@which $(MIGRATE_BIN) ||  ( curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz migrate && sudo mv migrate /usr/bin/migrate )

migrate-create: ## Create migration file with name
	migrate create -ext sql -dir internal/persistence/postgres/migrations -seq -digits 14 $(NAME)

migrate-up: check-migrate ## Run migrations
	$(MIGRATE) up

migrate-down: check-migrate ## Rollback migrations
	yes | $(MIGRATE) down

migrate-down-1: check-migrate ## Rollback last migration
	$(MIGRATE) down 1

dc-up:
	@docker-compose -f ./deployment/docker-compose.yaml up -d

dc-stop:
	@docker-compose -f ./deployment/docker-compose.yaml stop

dc-clean:
	@cd ./dev ; docker-compose stop ; docker-compose rm -f

dc-show:
	@docker container ls --format "{{.Names}} [{{.Ports}}]"

imports:
	goimports -w .

fmt:
	go fmt ./...

mock-repo:
	mockgen --build_flags=--mod=mod -destination=mocks/repo.go -package=mocks gitlab.yalantis.com/payments/payroll/service Repo

precomit: test imports fmt lint
