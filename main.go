package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/odise/moxy"
)

func readConfig(path string) (*Config, error) {

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	config, err := ParseConfig(string(buf))
	if err != nil {
		return nil, err
	}

	return config, nil
}

//AddSecurityHeaders comment
func AddSecurityHeaders(request *http.Request, response *http.Response) {
	response.Header.Del("X-Powered-By")
	response.Header.Set("X-Super-Secure", "Yes!!")
}

func main() {

	cfg, err := readConfig("./config")
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	for _, p := range cfg.Proxies {
		createProxies(&p, router)
	}

	app := negroni.New()
	app.UseHandler(router)
	app.Run(":3009")
}

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	log.Printf("revprox: This is an authenticated request")
	log.Printf("revprox: Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims {
		log.Printf("revprox: %s :\t%#v\n", k, v)
	}
})

func addHost(r *mux.Router, hostname string) *mux.Router {
	s := r.Host(hostname).Subrouter()
	return s
}

func addSecureRoute(router *mux.Router, proxy *moxy.ReverseProxy, route string, jwtMiddleware *jwtmiddleware.JWTMiddleware) {
	router.Handle(route, negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(myHandler),
		negroni.HandlerFunc(proxy.HandlerWithNext),
	))
}

func addPublicRoute(router *mux.Router, proxy *moxy.ReverseProxy, route string) {
	router.Handle(route, negroni.New(
		negroni.HandlerFunc(proxy.HandlerWithNext),
	))
}
