package endpoint

import (
	"net/http"
)

type Endpoint interface {
	Handler(res http.ResponseWriter, req *http.Request)
}
