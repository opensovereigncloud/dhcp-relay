FROM golang:1.23.0 AS builder
ARG OS
ARG ARCH

WORKDIR /workspace
# Copy go modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/
# Build
RUN CGO_ENABLED=0 GOOS="${OS:-linux}" GOARCH="${ARCH}" go build -a -o dhcp-relay cmd/main.go

FROM debian:bookworm-20240722-slim as installer
RUN apt update
# Install dhcrelay
RUN DEBIAN_FRONTEND=noninteractive apt install -y isc-dhcp-relay

FROM debian:bookworm-20240722-slim
WORKDIR /
COPY --from=builder /workspace/dhcp-relay .
COPY --from=installer /usr/sbin/dhcrelay /usr/sbin/dhcrelay
USER 65535:65535
ENTRYPOINT ["/dhcp-relay"]
CMD [""]