package authentication

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
)

type Authentification struct {
	allowedAccounts map[string]Account
}

func New(allowedAccounts map[string]Account) *Authentification {
	auth := new(Authentification)
	auth.allowedAccounts = allowedAccounts
	return auth
}

func (a *Authentification) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		appId := req.Header.Get(HEADER_APPID_KEY)
		log.Debugf("Authenticate AppId: %s", appId)
		if a.authenticate(appId, req.Header.Get(HEADER_APIKEY_KEY)) == false {
			log.Infof("Authentication for AppId '%s' failed!", appId)
			res.WriteHeader(http.StatusForbidden)
			msg := fmt.Sprintf(
				"Authetification failed! Ensure you send correct %s & %s within your header.",
				HEADER_APPID_KEY, HEADER_APIKEY_KEY,
			)
			res.Write([]byte(msg))
			return
		}
		log.Infof("Authentification for AppId '%s' passed", appId)
		next.ServeHTTP(res, req)
	})
}

func (a *Authentification) authenticate(appId string, apiKey string) bool {
	if appId != "" && a.allowedAccounts[appId].ApiKey == apiKey {
		return true
	}

	return false
}
