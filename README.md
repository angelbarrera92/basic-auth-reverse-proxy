# Basic Auth Reverse Proxy

This project offers a way to securize your backends with a basic golang reverse proxy with basic auth configuration.

![Logo](assets/logo-small.png)

## TLDR;

Expose your services with basic authentication.

```bash
$ wget -O basic-auth-reverse-proxy.tar.gz -q https://github.com/angelbarrera92/basic-auth-reverse-proxy/releases/download/v0.1.2/basic-auth-reverse-proxy_0.1.2_linux_amd64.tar.gz
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

## Configuration

The server can be configured via a file or environment variables.

### File

The configuration file is parsed as YAML input. The default location is a file
named *authn.yaml* in the current working directory.
The file contains authentication principals for the reverse proxy.

```yaml
---
users:
  - username: Angel
    password: Barrera
  - username: Pepe
    password: Gotera
```

### Environment variables

The environment variables containing configuration values all have the prefix
**BARP_**. Both daemon settings as well as authentication principals can be
provided this way:

* **BARP_UPSTREAM**

  The proxied URL

* **BARP_PORT**

  The port to bind to

* **BARP_REALM**

  The basic authentication realm identifier

Those environment variables take precedence over CLI arguments. Authentication
principals are additive (i.e. they will not remove entities from
the configuration file).
Principals are defined via environment variables prefixed with *BARP_USERNAME_*
(e.g. BARP_USERNAME_1=test). The values are the authentication username. The
password can be defined using the prefix *BARP_PASSWORD_*
(e.g. BARP_PASSWORD_1=password_for_test). If no password is defined a random
string is generated and emitted on stdout. The generated password is not very
secure and should only be used for testing/debugging purposes.
