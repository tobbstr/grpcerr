package grpcerr

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
)

func ExampleNewHttpResponseEncodeWriter() {
	w := httptest.NewRecorder()
	encodeAndWrite := NewHttpResponseEncodeWriter(w)

	unimplementedGRPCError := NewUnimplemented("")

	if err := encodeAndWrite(unimplementedGRPCError).AsJSON(); err != nil {
		panic(err)
	}

	result := w.Result()
	defer result.Body.Close()
	httpResponseBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nHTTP Status code: %d\n", result.StatusCode)
	fmt.Printf("HTTP Body:\n%s\n", string(httpResponseBody))
}
