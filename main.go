package main

import (
	"net/url"
	"net/http/httputil"
	"net/http"
	"github.com/angelbarrera92/basic-auth-reverse-proxy/proxy"
)


func main() {
	upstreamURL, _ := url.Parse("https://httpbin.org")
	reverseProxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	http.HandleFunc("/", proxy.ReverseProxyHandler(reverseProxy, upstreamURL))
	if err := http.ListenAndServe(":11811", nil); err != nil {
		panic(err)
	}
}
