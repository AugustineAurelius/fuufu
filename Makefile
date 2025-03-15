.PHONY: install generate import api air release migration

install:
	go install github.com/AugustineAurelius/eos@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/air-verse/air@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/goreleaser/goreleaser/v2@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

build-frontend:
	yarn  --cwd ./frontend run build

build: build-frontend
	go build  -o bin/

run:
	./bin/fuufu.exe serve
generate:
	go generate ./...
	make import

import:
	goimports -w .

api:
	oapi-codegen --config=api/todo/config.yaml api/todo/api.yaml

air:
	air -c .air.toml

air-serve:
	air -c .air.toml -build.bin=tmp\\main.exe serve -c=internal/config/config.yaml

release:
	goreleaser release

fmt:
	go fmt ./...

 migration:
	goose create $(name) go --dir=db/postgres/migrations