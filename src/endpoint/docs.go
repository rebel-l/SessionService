package endpoint

import (
	"net/http"
	log "github.com/sirupsen/logrus"
)

func InitDocsEndpoint()  {
	log.SetLevel(log.DebugLevel)
	log.Debug("Docs endpoint: Init ...")
	fs := http.FileServer(http.Dir("docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))
	log.Debug("Docs endpoint: initialized!")
}
