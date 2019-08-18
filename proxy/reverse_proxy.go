package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func modifyRequest(r *http.Request, upstreamURL *url.URL) {
	// Update the headers to allow for SSL redirection
	r.URL.Host = upstreamURL.Host
	r.URL.Scheme = upstreamURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Host)
	user, _, _ := r.BasicAuth()
	r.Header.Set("X-Scope-OrgID", user)
	r.Host = upstreamURL.Host
}

// ReverseProxyHandler Handle every proxt request
func ReverseProxyHandler(reverseProxy *httputil.ReverseProxy, upstreamURL *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		modifyRequest(r, upstreamURL)
		reverseProxy.ServeHTTP(w, r)
	}
}
