TAG=${shell git describe --tags}

.PHONY: make
make:
	CGO_ENABLED=0 GOOS=linux \
	go build \
		-buildmode=pie \
		-ldflags="-s -w \
			-X main.version=@${TAG}" \
		-o ./onetim3 ./main.go

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
	docker build --build-arg TAG=${TAG} --tag ghcr.io/cyllective/onetim3:${TAG} .
	docker build --build-arg TAG=${TAG} --tag ghcr.io/cyllective/onetim3:latest .

.PHONY: docker-push
docker-push:
	docker push ghcr.io/cyllective/onetim3:${TAG}
	docker push ghcr.io/cyllective/onetim3:latest