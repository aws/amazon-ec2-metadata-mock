FROM golang:1.17 as builder

## GOLANG env
ARG GOPROXY="https://proxy.golang.org|direct"
ARG GO111MODULE="on"

# Copy go.mod and download dependencies
WORKDIR /amazon-ec2-metadata-mock
COPY go.mod .
COPY go.sum .
RUN go mod download

## GOLANG env
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

# Build
COPY . .
RUN make build
# In case the target is build for testing:
# $ docker build  --target=builder -t test .
ENTRYPOINT ["/amazon-ec2-metadata-mock/build/ec2-metadata-mock"]

# Build the final image with only the binary
FROM scratch
WORKDIR /
COPY --from=builder /amazon-ec2-metadata-mock/build/ec2-metadata-mock .
COPY THIRD_PARTY_LICENSES .
ENTRYPOINT ["/ec2-metadata-mock"]
