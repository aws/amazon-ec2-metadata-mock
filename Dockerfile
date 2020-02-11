# Build the manager binary
FROM golang:1.14 as builder

## GOLANG env
ARG GOPROXY="direct"
ARG GO111MODULE="on"
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

# Copy go.mod and download dependencies
WORKDIR /amazon-ec2-metadata-mock
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN make build
# In case the target is build for testing:
# $ docker build  --target=builder -t test .
ENTRYPOINT ["/amazon-ec2-metadata-mock/build/amazon-ec2-metadata-mock"]

FROM scratch
WORKDIR /
COPY --from=builder /amazon-ec2-metadata-mock/build/amazon-ec2-metadata-mock .
COPY THIRD_PARTY_LICENSES .
ENTRYPOINT ["/amazon-ec2-metadata-mock"]
