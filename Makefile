APP="svachan"
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES = $(shell  go list ./... | grep -v "vendor")

DB_USER ?= postgres
DB_PASS ?= password
DB_HOST ?= localhost
DB_PORT ?= 5432

setup:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	GO111MODULE=off go get -u github.com/matm/gocov-html
	export GO111MODULE=on

copy-config:
	cp application.yml.sample application.yml

build: fmt vet lint compile

compile:
	GO111MODULE=on go mod vendor
	mkdir -p out/
	go build -o $(APP_EXECUTABLE)

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@for p in $(ALL_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

test:
	ENVIRONMENT=test go test -cover  ./... -coverprofile cover.out
	ENVIRONMENT=test go tool cover -func cover.out

test-cov:
	gocov test ${ALL_PACKAGES} > docs/cov.json

test-cover-html:
	@echo "\nEXPORTING RESULTS TO COVERAGE.HTML..."
	gocov-html docs/cov.json > docs/coverage.html
	@echo 'TEST RESULTS EXPORTED TO DOCS/COVERAGE.HTML'

test-cov-report:
	@echo "\nGENERATING TEST REPORT."
	gocov report docs/cov.json

migration-up:
	goose -dir db/migrations/ postgres "user=$(DB_USER) password=$(DB_PASS) dbname=$(DB_NAME) sslmode=disable" up

migration-down:
	goose -dir db/migrations/ postgres "user=$(DB_USER) password=$(DB_PASS) dbname=$(DB_NAME) sslmode=disable" down