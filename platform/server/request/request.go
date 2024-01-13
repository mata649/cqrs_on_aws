package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrRequestError = errors.New("Error binding request")

func Binding[T interface{}](req *T, r *http.Request) error {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return fmt.Errorf("%w: %s", ErrRequestError, body)
	}
	err = json.Unmarshal(body, req)
	if err != nil {

		return fmt.Errorf("%w: %s", ErrRequestError, body)
	}

	return nil
}
