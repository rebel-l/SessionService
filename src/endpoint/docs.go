package endpoint

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func InitDocsEndpoint(r *mux.Router)  {
	log.Debug("Docs endpoint: Init ...")

	fs := http.FileServer(http.Dir("docs"))
	r.Handle("/docs/", http.StripPrefix("/docs/", fs))

	log.Debug("Docs endpoint: initialized!")
}
