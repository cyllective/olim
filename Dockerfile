FROM golang:1.26.0-alpine AS build

RUN wget "https://github.com/upx/upx/releases/download/v5.0.2/upx-5.0.2-amd64_linux.tar.xz" -q -O /tmp/upx.tar.xz && \
	tar -C /tmp -xf /tmp/upx.tar.xz && \
	mv /tmp/upx-*/upx /tmp/upx

RUN apk update && apk add git make

WORKDIR /src
RUN mkdir -p /src/bin
RUN git config --global --add safe.directory /src

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=bind,source=go.sum,target=go.sum \
	--mount=type=bind,source=go.mod,target=go.mod \
	go mod download

COPY . .
ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod/ \
	--mount=type=cache,target="/root/.cache/go-build" \
	make

RUN /tmp/upx --no-color -q --best -o /src/bin/olim /src/olim

FROM alpine:latest

COPY --from=build /src/bin/olim /olim

EXPOSE 8080

ENTRYPOINT ["/olim"]