package sessionservice

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/rebel-l/Session/1.0.0/",
		Index,
	},

	Route{
		"DeleteSession",
		"DELETE",
		"/rebel-l/Session/1.0.0/session",
		DeleteSession,
	},

	Route{
		"LoadSession",
		"GET",
		"/rebel-l/Session/1.0.0/session",
		LoadSession,
	},

	Route{
		"UpdateSession",
		"PUT",
		"/rebel-l/Session/1.0.0/session",
		UpdateSession,
	},

}