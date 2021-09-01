# Build the manager binary
FROM golang:1.17 as mod

## GOLANG env
ARG GOPROXY="direct"
ARG GO111MODULE="on"
ARG GOPATH=/go
ARG GOCACHE=/go

# Copy go.mod and download dependencies
WORKDIR /amazon-ec2-metadata-mock
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the manager binary
FROM golang:1.17 as builder

## GOLANG env
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
ARG GOPATH=/go
ARG GOCACHE=/go

# Copy go.mod and download dependencies
WORKDIR /amazon-ec2-metadata-mock
COPY --from=mod $GOCACHE $GOCACHE
COPY --from=mod $GOPATH/pkg/mod $GOPATH/pkg/mod

# Build
COPY . .
RUN make build
# In case the target is build for testing:
# $ docker build  --target=builder -t test .
ENTRYPOINT ["/amazon-ec2-metadata-mock/build/ec2-metadata-mock"]

FROM scratch
WORKDIR /
COPY --from=builder /amazon-ec2-metadata-mock/build/ec2-metadata-mock .
COPY THIRD_PARTY_LICENSES .
ENTRYPOINT ["/ec2-metadata-mock"]
