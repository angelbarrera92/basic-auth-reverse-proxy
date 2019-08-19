package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/angelbarrera92/basic-auth-reverse-proxy/proxy"
	"gopkg.in/urfave/cli.v1"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func serve(c *cli.Context) error {
	upstream := c.String("upstream")
	port := c.Int("port")
	authConfigPath := c.String("auth-config")
	realm := c.String("realm")

	authConfig, err := proxy.ParseConfig(&authConfigPath)

	if err != nil {
		log.Fatalf("Can not read auth configuration file: %v", err)
		return err
	}

	upstreamURL, _ := url.Parse(upstream)
	reverseProxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	http.HandleFunc("/", proxy.BasicAuth(proxy.ReverseProxyHandler(reverseProxy, upstreamURL, authConfig), *authConfig, realm))
	serveAt := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(serveAt, logRequest(http.DefaultServeMux)); err != nil {
		log.Fatalf("Reverse Proxy can not start %v", err)
		return err
	}

	return nil
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	app := cli.NewApp()
	app.Name = "Basic Auth Reverse Proxy"
	app.Usage = "Makes your upstream service secure"
	app.Version = version
	app.Author = "Ángel Barrera - @angelbarrera92"
	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "Runs the reverse proxy",
			Action: serve,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Usage: "Port used to expose this reverse proxy",
					Value: 11811,
				}, cli.StringFlag{
					Name:  "upstream",
					Usage: "Upstream server. Server that will be protected by this reverse proxy",
					Value: "https://httpbin.org",
				}, cli.StringFlag{
					Name:  "realm",
					Usage: "Reverse proxy realm",
					Value: "My Reverse Proxy",
				}, cli.StringFlag{
					Name:  "auth-config",
					Usage: "AuthN yaml configuration file path",
					Value: "authn.yaml",
				},
			},
		},
	}
	app.Run(os.Args)
}
