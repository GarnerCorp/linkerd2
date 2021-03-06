ARG BUILDPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM docker.io/curlimages/curl:7.75.0  as await
WORKDIR /tmp
ARG LINKERD_AWAIT_VERSION=v0.2.1
ARG TARGETARCH
RUN curl -fsSvLo linkerd-await https://github.com/olix0r/linkerd-await/releases/download/release%2F${LINKERD_AWAIT_VERSION}/linkerd-await-${LINKERD_AWAIT_VERSION}-${TARGETARCH} && chmod +x linkerd-await

# Precompile key slow-to-build dependencies
FROM --platform=$BUILDPLATFORM golang:1.16.2-alpine as go-deps
WORKDIR /linkerd-build
COPY go.mod go.sum ./
COPY bin/install-deps bin/
RUN go mod download
ARG TARGETARCH
RUN ./bin/install-deps $TARGETARCH

## compile metrics-apiservice
FROM go-deps as golang
WORKDIR /linkerd-build
COPY pkg pkg
COPY controller controller
COPY viz/metrics-api viz/metrics-api
COPY viz/pkg viz/pkg

ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /out/metrics-api -tags prod -mod=readonly -ldflags "-s -w" ./viz/metrics-api/cmd

## package runtime
FROM scratch
COPY --from=await /tmp/linkerd-await /
COPY LICENSE /linkerd/LICENSE
COPY --from=golang /out/metrics-api /metrics-api
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/linkerd-await", "--", "/metrics-api"]
