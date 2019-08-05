package proxy

import (
	"crypto/subtle"
	"net/http"
)

// BasicAuth Protects a handler with a users authn backend
func BasicAuth(handler http.HandlerFunc, users Authn, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok {
			writeUnauthorisedResponse(w, realm)
			return
		}

		authn := false
		for _, v := range users.Users {
			if subtle.ConstantTimeCompare([]byte(user), []byte(v.Username)) == 1 && subtle.ConstantTimeCompare([]byte(pass), []byte(v.Password)) == 1 {
				authn = true
				break
			}
		}

		if !authn {
			writeUnauthorisedResponse(w, realm)
			return
		}

		handler(w, r)
	}
}

func writeUnauthorisedResponse(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised\n"))
}
