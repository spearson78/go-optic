package ohttp

import (
	"fmt"
	"net/http"
)

type HttpError int

func (h HttpError) Error() string {
	return fmt.Sprintf("%v %v", int(h), http.StatusText(int(h)))
}
