package proxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func modifyRequest(r *http.Request, upstreamURL *url.URL, users *Authn) {
	// Update the headers to allow for SSL redirection
	r.URL.Host = upstreamURL.Host
	r.URL.Scheme = upstreamURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Host)
	userName, _, _ := r.BasicAuth()
	orgID, _ := findOrgIDforUser(userName, users)
	r.Header.Set("X-Scope-OrgID", orgID)
	r.Host = upstreamURL.Host
}

// ReverseProxyHandler Handle every proxt request
func ReverseProxyHandler(reverseProxy *httputil.ReverseProxy, upstreamURL *url.URL, users *Authn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		modifyRequest(r, upstreamURL, users)
		reverseProxy.ServeHTTP(w, r)
	}
}

func findOrgIDforUser(userName string, users *Authn) (string, error) {
	for _, v := range users.Users {
		if v.Username == userName {
			return v.OrgID, nil
		}
	}
	return "", errors.New("User not found")
}
