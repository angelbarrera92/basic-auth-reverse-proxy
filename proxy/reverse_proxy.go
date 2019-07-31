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
	r.Host = upstreamURL.Host
}

func ReverseProxyHandler(reverseProxy *httputil.ReverseProxy, upstreamURL *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		modifyRequest(r, upstreamURL)
		reverseProxy.ServeHTTP(w, r)
	}
}
