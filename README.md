# TL;DR

This library enables API developers to use the gRPC error model for both gRPC and HTTP.
It provides an easy to use API for instantiating gRPC errors such as Unavailable, Resource Exhausted, Invalid argument etc,
in addition to easily adding metadata such as a request ID, stack traces and more. It also facilitates encoding and writing of
gRPC errors to a http.ResponseWriter. As of now, it only supports JSON-encoding.

# Getting Started

So you want to get started using this library? Follow the instructions in this section.

## Installation

Given the project that wants to use this library is using Go Modules, installing it is as easy as entering the following command:

```
go get github.com/tobbstr/grpcerr
```

This'll add it to the list of dependencies in the go.mod file.

# Usage

## Instantiation of a gRPC error

```go
import (
    "github.com/tobbstr/grpcerr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func main() {
    // Error details
    errInfo := &errdetails.ErrorInfo{
        Reason: "Token expired",
        Domain: "Authorization",
        Metadata: map[string]string{
            "TokenExpired": "2006-01-02 15:04:05",
        },
    }

    // Instantiation of the gRPC error
    permissionDenied, err := grpcerr.NewPermissionDenied("", errInfo)
    if err != nil {
        // handle error
    }
}
```

## Converting gRPC errors to error

```go
import (
    "github.com/tobbstr/grpcerr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func main() {
    // Converting a gRPC error to an error
    err := grpcerr.ToError(permissionDenied)
}
```

## Converting errors to gRPC errors

```go
import (
    "github.com/tobbstr/grpcerr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func main() {
    // Converting an error to a gRPC error
    gRPCErr := grpcerr.FromError(err)
}
```

## Returning gRPC error from an HTTP API

```go
import (
    "net/http"
    "github.com/tobbstr/grpcerr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// controller definition omitted

func (c *controller) awesomeEndpoint(w http.ResponseWriter, r *http.Request) {
    // ... Do bunch of stuff ...

    // Handle errors returned from an invocation
    err = vitalFunc(...)
    if err != nil {
        encodeAndWrite := grpcerr.NewHttpResponseEncodeWriter(w)
        if err = encodeAndWrite(err).AsJSON(); err != nil {
            // Log and handle error. The error should normally be nil.
        }
        return
    }

    // Continue normal processing...
}
```

# Roadmap

See the [open issues](https://github.com/tobbstr/grpcerr/issues) for a list of proposed features (and known issues).

# Contributing

Contributing can be done in more than one way. One way is to hit the star button. Another is to add to or improve the existing repo content. The following steps guide you how to do that. In any case, contributions are greatly appreciated.

1. Clone this repo by entering this command in a terminal
    ```sh
    git clone https://github.com/tobbstr/grpcerr.git
    ```
2. Navigate to the repo folder
3. Create a feature branch
    ```sh
    git checkout -b feature-name
    ```
4. Make the changes you want and commit them
5. Push the branch to GitHub
    ```sh
    git push origin feature-name
    ```
6. Open a pull request

# License

Distributed under the [MIT](LICENSE) license.