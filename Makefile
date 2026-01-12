.PHONY: make
make:
	CGO_ENABLED=0 GOOS=linux \
	go build \
		-buildmode=pie \
		-ldflags="-s -w \
			-X main.version=@`git describe --tags --abbrev=0`" \
		-o ./olim ./main.go

.PHONY: lint
lint:
	deadcode ./...
	modernize ./...
	goimports-reviser -format ./...
	golangci-lint run

.PHONY: updatepackages
updatepackages:
	go mod tidy
	go get ${shell go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all}

.PHONY: docker-image
docker-image:
	docker build --tag ghcr.io/cyllective/olim:${shell git describe --tags --abbrev=0} .
	docker build --tag ghcr.io/cyllective/olim:latest .

.PHONY: docker-push
docker-push:
	docker push ghcr.io/cyllective/olim:${shell git describe --tags --abbrev=0}
	docker push ghcr.io/cyllective/olim:latest