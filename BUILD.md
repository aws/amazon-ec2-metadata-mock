# Amazon EC2 Metadata Mock: Build Instructions

## Install Go version 1.14+

There are several options for installing go:

1. If you're on mac, you can simply `brew install go`
2. If you'd like a flexible go installation manager consider using gvm https://github.com/moovweb/gvm
3. For all other situations use the official go getting started guide: https://golang.org/doc/install

## Compile

This project uses `make` to organize compilation, build, and test targets.

To compile cmd/amazon-ec2-metadata-mock.go which will build the full static binary and pull in dependent packages:
```
make build
```

The resulting binary will be in the generated `build/` dir

```
$ make build

$ ls build/
ec2-metadata-mock
```

## Test

You can execute the unit tests for AEMM with `make`:

```
make unit-test
```


### Run All Tests

The full suite includes unit tests, integration tests, and more. See the full list in the [makefile](https://github.com/aws/amazon-ec2-metadata-mock/blob/main/Makefile): 

```
make test
```

## Format

To keep our code readable with go conventions, we use `goimports` to format the source code.
Make sure to run `goimports` before you submit a PR or you'll be caught by our tests! 

You can use the `make fmt` target as a convenience
```
make fmt
```
