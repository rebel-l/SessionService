package authentication

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
)

type Authentification struct {
	allowedAccounts []Account
	//Handler func(next http.Handler) http.Handler
}

//type Middleware interface {
//	handler(next http.Handler) http.Handler
//}

func New(allowedAccounts []Account) *Authentification {
	auth := new(Authentification)
	auth.allowedAccounts = allowedAccounts
	return auth
}

func (a *Authentification) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		clientAccount := Account{
			AppId:  req.Header.Get(HEADER_APPID_KEY),
			ApiKey: req.Header.Get(HEADER_APIKEY_KEY),
		}
		log.Debugf("Authenticate AppId: %s", clientAccount.AppId)
		if a.authenticate(clientAccount) == false {
			log.Infof("Authentication for AppId '%s' failed!", clientAccount.AppId)
			res.WriteHeader(http.StatusForbidden)
			msg := fmt.Sprintf(
				"Authetification failed! Ensure you send correct %s & %s within your header.",
				HEADER_APPID_KEY, HEADER_APIKEY_KEY,
			)
			res.Write([]byte(msg))
			return
		}
		log.Infof("Authentification for AppId '%s' passed", clientAccount.AppId)
		next.ServeHTTP(res, req)
	})
}

func (a *Authentification) authenticate(clientAccount Account) bool {
	valid := false
	for _, account := range a.allowedAccounts {
		if account.AppId != clientAccount.AppId {
			continue
		}

		if account.ApiKey == clientAccount.ApiKey {
			valid = true
		}
	}
	return valid
}
