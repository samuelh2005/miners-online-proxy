FROM --platform=$BUILDPLATFORM golang:1.26.2 AS build

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.sum ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY plugins ./plugins
COPY gate.go ./

# Automatically provided by the buildkit
ARG TARGETOS TARGETARCH

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -a -o gate gate.go

# Move binary into final image
FROM --platform=$BUILDPLATFORM eclipse-temurin:25.0.1_8-jre-alpine AS app
COPY --from=build /workspace/gate /
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh
#COPY config.yml /
ENTRYPOINT ["/entrypoint.sh"]