.PHONY: install generate import api air

install:
	go install github.com/AugustineAurelius/eos@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/air-verse/air@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

generate:
	go generate ./...
	make import


import:
	goimports -w .

api:
	oapi-codegen --config=api/config.yaml api/api.yaml

air:
	air -c .air.toml