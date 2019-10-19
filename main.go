package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/angelbarrera92/basic-auth-reverse-proxy/proxy"
	"gopkg.in/urfave/cli.v1"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func parseServeEnv(c *cli.Context) (upstream string, port int, realm string, err error) {
	upstream = c.String("upstream")
	port = c.Int("port")
	realm = c.String("realm")
	err = nil

	if env := os.Getenv("BARP_PORT"); len(env) != 0 {
		port, err = strconv.Atoi(env)
		if err != nil {
			err = fmt.Errorf("invalid environment variable BARP_PORT: %v", err)
			return
		}
	}

	if env := os.Getenv("BARP_UPSTREAM"); len(env) != 0 {
		upstream = env
	}

	if env := os.Getenv("BARP_REALM"); len(env) != 0 {
		realm = env
	}

	return
}

func serve(c *cli.Context) error {
	authConfigPath := c.String("auth-config")
	upstream, port, realm, err := parseServeEnv(c)

	if err != nil {
		log.Fatalf("Setup configuration failed: %v", err)
		return err
	}

	authConfig := proxy.NewAuthn()

	if err = authConfig.ParseFile(authConfigPath); err != nil {
		log.Fatalf("Can not read auth configuration file: %v", err)
		return err
	}

	if err = authConfig.ParseEnvironment(); err != nil {
		log.Fatalf("Unable to parse environment variables: %v", err)
		return err
	}

	if err = authConfig.Validate(); err != nil {
		log.Fatalf("Setup validation failed: %v", err)
		return err
	} else {
		log.Printf("Exposing %s:%d to %d users", upstream, port, len(authConfig.Users))
	}

	upstreamURL, _ := url.Parse(upstream)
	reverseProxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	http.HandleFunc("/", proxy.BasicAuth(proxy.ReverseProxyHandler(reverseProxy, upstreamURL), *authConfig, realm))
	serveAt := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(serveAt, nil); err != nil {
		log.Fatalf("Reverse Proxy can not start %v", err)
		return err
	}

	return nil
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
