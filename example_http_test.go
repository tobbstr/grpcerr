package grpcerr

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
)

func ExampleHttpResponseEncodeWriter() {
	w := httptest.NewRecorder()
	encodeAndWrite := NewHttpResponseEncodeWriter(w, nil)

	unimplementedGRPCError, err := NewUnimplemented("")
	if err != nil {
		panic(err)
	}

	if err = encodeAndWrite(unimplementedGRPCError).AsJSON(); err != nil {
		panic(err)
	}

	result := w.Result()
	defer result.Body.Close()
	httpResponseBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nStatus code: %d\n", result.StatusCode)
	fmt.Printf("Body:\n%s\n", string(httpResponseBody))
}
