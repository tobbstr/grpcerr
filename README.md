# TL;DR

This library enables API developers to use the gRPC error model for both gRPC and HTTP.
It provides an easy to use API for instantiating gRPC errors such as Unavailable, Resource Exhausted, Invalid argument etc,
in addition to easily adding metadata such as a request ID, stack traces and more. It also facilitates encoding and writing of
gRPC errors to an http.ResponseWriter. As of now, it only supports JSON-encoding.

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
    errInfo := &grpcerr.ErrorInfo{
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

    err = vitalFunc(...)
    // Handles errors returned from an invocation.
    // Note! err must be a gRPC error instantiated using this library
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

This is an example of a HTTP Response Body for an InvalidArgument gRPC error:

```json
{
    "code":3,
    "message":"dummy-msg",
    "details":[
        {
            "@type":"type.googleapis.com/google.rpc.BadRequest",
            "fieldViolations":[
                {
                    "field":"dummy-field-violation-field",
                    "description":"dummy-field-violation-desc"
                }
            ]
        }
    ]
}
```

In terms of HTTP headers the HTTP status code is translated from the gRPC error code, but is overridable using
options when calling `grpcerr.NewHttpRresponseEncodeWriter(w, opts...)`. The below is an example of this:

```go
// defining the override
withStatusOK := func(w http.ResponseWriter) {
    w.WriteHeader(http.StatusOK)
}

// use the option like this
encodeAndWrite := NewHttpResponseEncodeWriter(w, withStatusOK)
```

## Wrapping of errors are supported

```go
// Somewhere in a service far, far away...
func (svc *service) OrchestrateImportantAction() error {
    // Something goes awry
    if err != nil {
        violations := []grpcerr.PreconditionFailure{
            {
                Type: "Account missing",
                Subject: "john.doe@example.com",
                Description: "Callers must supply account information in request",
            },
        }

        failedPrecondition, err := grpcerr.NewFailedPrecondition("", violations)
        if err != nil {
            // handle error
        }

        // Note! The returned error is a wrapped one, having the gRPC error as its root error. It does not matter how many times the gRPC error gets wrapped, as long as it's the root error.
        return fmt.Errorf("could not perform something cool: %w", failedPrecondition)
    }
}

// Meanwhile in an HTTP controller not that far away
func (c *controller) performImportantAction(w http.ResponseWriter, r *http.Request) {
    // invoke svc, which returns a wrapped gRPC error instantiated using this library
    err := c.svc.OrchestrateImportantAction()
    if err != nil {
        // ... Log the error ...

        // write an HTTP response from the root error which is the gRPC error
        encodeAndWrite := grpcerr.NewHttpResponseEncodeWriter(w)
        if err = encodeAndWrite(err).AsJSON(); err != nil {
            // handle error
        }
        return
    }
}
```

## Using gRPC errors in gRPC APIs

```go
func (c *gRPCController) AwesomeEndpoint(ctx context.Context, req *AwesomeRequest) (*AwesomeResponse, error) {
    err := c.svc.OrchestrateImportantAction()
    if err != nil {
        // ... Log the error ...

        // no additional processing required, return the error as is
        return nil, err
    }
}
```

## Using gRPC errors in gRPC clients

```go
response, err := client.AwesomeEndpoint(...)
if err != nil {
    // gets the gRPC code
    code := grpcerr.Code(err)

    // gets the gRPC message
    message := grpcerr.Message(err)

    // gets the RequestInfo from the gRPC error
    requestInfo := grpcerr.RequestInfoFrom(err)

    // gets the DebugInfo from the gRPC error
    debugInfo := grpcerr.DebugInfoFrom(err)

    // gets the HelpLinks from the gRPC error
    helpLinks := grpcerr.HelpLinksFrom(err)

    // gets the LocalizedMessage from the gRPC error
    localizedMessage := grpcerr.LocalizedMessageFrom(err)

    // gets the QuotaViolations from the gRPC error
    quotaViolations := grpcerr.QuotaViolationsFrom(err)

    // gets the FieldViolations from the gRPC error
    fieldViolations := grpcerr.FieldViolationsFrom(err)

    // etc ...
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