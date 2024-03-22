# Basic Auth Reverse Proxy

This project offers a way to securize your backends with a basic golang reverse proxy with basic auth configuration.

![Logo](assets/logo-small.png)

## TLDR;

Expose your services with basic authentication.

```bash
$ wget -O basic-auth-reverse-proxy.tar.gz -q https://github.com/angelbarrera92/basic-auth-reverse-proxy/releases/download/v0.1.6/basic-auth-reverse-proxy_0.1.6_linux_amd64.tar.gz
$ tar -zxvf basic-auth-reverse-proxy.tar.gz
README.md
basic-auth-reverse-proxy
$ cat >> authn.yaml <<EOL
users:
  - username: Angel
    password: Barrera
EOL
$ ./basic-auth-reverse-proxy serve
```

Then a local server is started. Try to access it:

```bash
$ curl http://localhost:11811/get
Unauthorised
$ curl http://Angel:Barrera@localhost:11811/get
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "Authorization": "Basic QW5nZWw6QmFycmVyYQ==",
    "Host": "httpbin.org",
    "User-Agent": "curl/7.58.0",
    "X-Forwarded-Host": "localhost:11811"
  },
  "origin": "127.0.0.1, 80.25.227.133, 127.0.0.1",
  "url": "https://localhost:11811/get"
}
```
