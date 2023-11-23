FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o lavalink-bot github.com/lavalink-devs/lavalink-bot

FROM alpine

COPY --from=build /build/lavalink-bot /bin/lavalink-bot

ENTRYPOINT ["/bin/lavalink-bot"]

CMD ["-config", "/var/lib/lavalink-bot.yml"]
