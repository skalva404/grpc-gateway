# grpc-gateway

This helps you provide your APIs in both gRPC and RESTful style at the same time.

![architecture introduction diagram](https://docs.google.com/drawings/d/12hp4CPqrNPFhattL_cIoJptFvlAqm5wLQ0ggqI5mkCg/pub?w=749&amp;h=370)

[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=https://github.com/grpc-gateway.git)

All the boilerplate you need to get started with writing grpc-gateway powered
REST services in Go.

## Running

Running `bootstrap.go` starts a web server on https://0.0.0.0:11000/. 
An OpenAPI Swagger UI is served on https://0.0.0.0:11000/.

## Requirements

Generating the files requires the `protoc` protobuf compiler.
Please install it according to the
[installation instructions](https://github.com/google/protobuf#protocol-compiler-installation)
for your specific platform.

## Getting started

After cloning the repo, there are a couple of initial steps;

1. Install the generate dependencies with `make install`.
   This will install `protoc-gen-go`, `protoc-gen-grpc-gateway`, `protoc-gen-swagger` and `statik` which
   are necessary for us to generate the Go, swagger and static files.
2. Finally, generate the files with `make generate`.
   If you encounter an error here, make sure you've installed
   `protoc` and it is accessible in your `$PATH`, and make sure
   you've performed step 1.

Now you can run the web server with `go run bootstrap.go`.
