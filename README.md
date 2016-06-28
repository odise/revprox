# revprox

This is a simple example on how to implement reverse proxy functionality coupled with JWT token checking in Go.

## Disclaimer

This code is far away from being complete. It wasn't written as part of API Gateway evaluation and might be extended (or not). The tool should run behind a load-balancer which is responsible for TLS termination and scaling. Feel welcome to provide feedback.

## Configuration

Configuration is done with the help of [HashiCorp Configuration Language](https://github.com/hashicorp/hcl) (HCL).

Each backend can be defined in a section like this:

```
proxy "service.example.com" {
  # JWT HS256 algorithm secret
  secret = "secret"
  # URL paths that need JWT tokens
  authpath = [ "/admin", "/" ]
  # publicly accessible URL paths
  publicpath = [ "/article/{category}" ]
  # upstreams
  target = [ "https://google.com", "http://google.de" ]
}
```

### Path definition

`publicpath` and `authpath` can be defined in the format `/path/{name}` or `/path/{name:pattern}`. Please have a look at [mux](https://github.com/gorilla/mux) for further documentation.

### JWT header

Secured paths need [JSON Web Tokens](https://jwt.io/#debugger-io) for request authentication. Here is a curl example:

```
curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6IjEyMzQ1Njc4OTAifQ.FMRJjb3KhS7PHH4Uwg6N06Htp1XD1vCC8y84zNK4WYA" -H 'Host: service.example.com' localhost:3009/article/gateway
```

Requests missing `Authorization: Bearer ...` header will be responded with HTTP error code `302` along with `Location: https://login.example.com?redirect_uri=` header.

### Upstream `target`

It is possible to define multiple backends randomly chosen on request time. Check [moxy/utils.go](https://github.com/odise/moxy/blob/master/utils.go) for further information.
