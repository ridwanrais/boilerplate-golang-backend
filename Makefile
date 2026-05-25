.PHONY: run generate tidy migrate atlas-install apply test test-integration

run:
	go run cmd/api/main.go

generate:
	go generate ./ent

tidy:
	go mod tidy

test:
	go test -short ./...

test-integration:
	go test ./... -v

migrate:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate name=<migration_name>"; \
		exit 1; \
	fi
	mkdir -p ent/migrate/migrations
	./bin/atlas migrate diff $(name) \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema" \
		--dev-url "docker://postgres/15/dev?search_path=public"

atlas-install:
	mkdir -p bin
	curl -sSfL https://release.ariga.io/atlas/atlas-darwin-arm64-latest -o bin/atlas
	chmod +x bin/atlas

apply:
	./bin/atlas migrate apply \
		--dir "file://ent/migrate/migrations" \
		--url "postgres://devuser:devpassword@localhost:5439/golang_boilerplate_db?sslmode=disable"
